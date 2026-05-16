package main

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/crossplane/function-sdk-go/resource/composed"
	"github.com/crossplane/function-sdk-go/response"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/crossplane/function-example/input/v1beta1"
)

// Function returns whatever response you ask it to.
type Function struct {
	fnv1beta1.UnimplementedFunctionRunnerServiceServer
	log logging.Logger
}

// RunFunction runs the Function.
func (f *Function) RunFunction(_ context.Context, req *fnv1beta1.RunFunctionRequest) (*fnv1beta1.RunFunctionResponse, error) {
	f.log.Info("Running GRAS logic function", "tag", req.GetMeta().GetTag())

	rsp := response.To(req, response.DefaultTTL)

	// Get the input to the function
	in := &v1beta1.Input{}
	if err := request.GetInput(req, in); err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get function input"))
		return rsp, nil
	}

	// Get the observed composite resource (XR)
	xr, err := request.GetObservedCompositeResource(req)
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get observed composite resource"))
		return rsp, nil
	}

	// Extract spec from XR
	spec, ok := xr.Resource.Object["spec"].(map[string]interface{})
	if !ok {
		f.log.Info("XR has no spec, skipping")
		return rsp, nil
	}

	// Define defaults from GRAS spec
	parentClaimRef, _ := spec["claimRef"].(map[string]interface{})
	parentClaimName, _ := parentClaimRef["name"].(string)
	parentNamespace, _ := parentClaimRef["namespace"].(string)

	asname, _ := spec["asname"].(string)
	if asname == "" {
		asname = parentClaimName
	}

	grasDefaults := map[string]interface{}{
		"asname":        asname,
		"grasversion":   spec["grasversion"],
		"clusterdomain": spec["clusterdomain"],
		"customheader":  spec["customheader"],
		"dev":           spec["dev"],
		"ingress":       spec["ingress"],
		"ssl":           spec["ssl"],
		"sslissuer":     spec["sslissuer"],
		"autoscaling":   spec["autoscaling"],
		"resources":     spec["resources"],
	}

	// Get existing desired resources to avoid overwriting others (like the license check step)
	desired, err := request.GetDesiredComposedResources(req)
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get desired composed resources"))
		return rsp, nil
	}

	// Process Grapis
	// Map to store grapi names for gruim mapping
	grapiMap := make(map[string]string)

	if grapis, ok := spec["grapis"].([]interface{}); ok {
		for i, g := range grapis {
			grapi, ok := g.(map[string]interface{})
			if !ok {
				continue
			}

			name, _ := grapi["name"].(string)
			grapiSpec, _ := grapi["spec"].(map[string]interface{})
			if grapiSpec == nil {
				grapiSpec = make(map[string]interface{})
			}

			// Apply defaults
			applyDefaults(grapiSpec, grasDefaults)

			// Propagate claimRef for patches in child compositions
			childName := name
			if in.ChildNameScheme != "" {
				res, err := f.renderTemplate(in.ChildNameScheme, map[string]string{
					"parent": parentClaimName,
					"child":  name,
				})
				if err == nil {
					childName = res
				}
			}
			grapiMap[name] = childName // Store for gruims

			grapiSpec["claimRef"] = map[string]interface{}{
				"apiVersion": "grsf.grpl.io/v1alpha1",
				"kind":       "GrappleApi",
				"namespace":  parentNamespace,
				"name":       childName,
			}

			// Create CompositeGrappleApi resource
			res := &composed.Unstructured{
				Unstructured: unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "grsf.grpl.io/v1alpha1",
						"kind":       "CompositeGrappleApi",
						"metadata": map[string]interface{}{
							"name": childName,
						},
						"spec": grapiSpec,
					},
				},
			}

			resourceName := fmt.Sprintf("grapi-%d", i)
			desired[resource.Name(resourceName)] = &resource.DesiredComposed{
				Resource: res,
			}
		}
	}

	// Process Gruims
	if gruims, ok := spec["gruims"].([]interface{}); ok {
		for i, g := range gruims {
			gruim, ok := g.(map[string]interface{})
			if !ok {
				continue
			}

			name, _ := gruim["name"].(string)
			gruimSpec, _ := gruim["spec"].(map[string]interface{})
			if gruimSpec == nil {
				gruimSpec = make(map[string]interface{})
			}

			// Apply defaults
			applyDefaults(gruimSpec, grasDefaults)

			// Propagate claimRef for patches in child compositions
			childName := name
			if in.ChildNameScheme != "" {
				res, err := f.renderTemplate(in.ChildNameScheme, map[string]string{
					"parent": parentClaimName,
					"child":  name,
				})
				if err == nil {
					childName = res
				}
			}
			gruimSpec["mapi"] = ""
			if in.MapiNamingScheme != "" {
				res, err := f.renderTemplate(in.MapiNamingScheme, map[string]string{
					"asname": asname,
				})
				if err == nil {
					gruimSpec["mapi"] = res
				}
			}
			gruimSpec["claimRef"] = map[string]interface{}{
				"apiVersion": "grsf.grpl.io/v1alpha1",
				"kind":       "GrappleUiModule",
				"namespace":  parentNamespace,
				"name":       childName,
			}

			// Create CompositeGrappleUiModule resource
			res := &composed.Unstructured{
				Unstructured: unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "grsf.grpl.io/v1alpha1",
						"kind":       "CompositeGrappleUiModule",
						"metadata": map[string]interface{}{
							"name": childName,
						},
						"spec": gruimSpec,
					},
				},
			}

			resourceName := fmt.Sprintf("gruim-%d", i)
			desired[resource.Name(resourceName)] = &resource.DesiredComposed{
				Resource: res,
			}
		}
	}

	// Set desired resources in response
	if err := response.SetDesiredComposedResources(rsp, desired); err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot set desired composed resources"))
		return rsp, nil
	}

	return rsp, nil
}

func applyDefaults(childSpec map[string]interface{}, defaults map[string]interface{}) {
	for k, v := range defaults {
		if v == nil {
			continue
		}
		// If field is missing or null in child, use default from parent
		if val, exists := childSpec[k]; !exists || val == nil {
			childSpec[k] = v
		}
	}
}

// renderTemplate renders a string template with the given data
func (f *Function) renderTemplate(tmplStr string, data interface{}) (string, error) {
	if tmplStr == "" {
		return "", nil
	}
	tmpl, err := template.New("tmpl").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
