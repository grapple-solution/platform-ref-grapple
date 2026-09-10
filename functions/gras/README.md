# Testing a Sample/Simple Crossplane Function


## Prerequisites

Ensure you have the following:
- Docker installed
- Go installed
- Access to a Kubernetes cluster

## Steps to Test the Crossplane Function

1. **Pull the Latest Changes:**
   - Clone the repository or navigate to your local `grsf-license-function` repository.
   - Pull the latest changes:
     ```bash
     git pull origin main
     ```

2. **Navigate to the License Function Directory:**
   - Move to the `grsf-license-function` directory:
     ```bash
     cd grsf-license-function
     ```

3. **Build the Docker Image:**
   - Open a terminal and run the following command to build the Docker image:
     ```bash
     docker buildx build --platform "linux/amd64,linux/arm64" . --tag=runtime
     ```
     Note: the process is not working for me, when building multi-platform image on my arm64 based machine.
           I have to build the image for amd64 and then download it as a tar and use the tar to build the package
     ```bash
     docker buildx build --platform "linux/amd64" . --tag=runtime
     docker save runtime -o runtime.tar
     ```


4. **Run the Function:**
   - Once the Docker image is built, run the function:
     ```bash
     go run . --insecure
     ```
   - The function will now be available on port `9443`.

5. **Testing Different Use Cases:**

   - **Add/Update the `LIC` Secret:**
     - Use the following commands to encode the `LIC` value and update the `grsf-config` secret in the `grpl-system` namespace:
       ```bash
       echo -n "<LIC value>" | base64
       ```
     - Copy the base64-encoded output.
     - Update the secret:
       ```bash
       kubectl patch secret grsf-config -n grpl-system --type='json' -p='[{"op": "add", "path": "/data/LIC", "value": "<copied_value>"}]'
       ```

   - **Edit the `composition.yaml` File:**
     - Update the `grsf-license-function/example/composition.yaml` file to set the email parameter’s value as required for the test case.

   - **Render and Test:**
     - Open a new terminal, navigate to the `example` directory, and run:
       ```bash
       cd grsf-license-function/example
       crossplane beta render xr.yaml composition.yaml functions.yaml
       ```
     - If successful, it will produce the following output:
       ```yaml
       ---
       apiVersion: example.crossplane.io/v1
       kind: XR
       metadata:
         name: example-xr
       status:
         conditions:
         - lastTransitionTime: "2024-01-01T00:00:00Z"
           reason: Available
           status: "True"
           type: Ready
       ```

   - **Verify the `GRAPPLE_LICENSE` Secret:**
     - To confirm that the `GRAPPLE_LICENSE` secret was added or updated, run:
       ```bash
       kubectl get secret -n grpl-system grsf-config -o jsonpath="{.data.GRAPPLE_LICENSE}" | base64 --decode
       ```

6. **Repeat Testing:**
   - You can repeat step 5 to test different use cases as needed.



## Steps to Release the Crossplane Function


1. **Pull the Latest Changes:**
   - Clone the repository or navigate to your local `grsf-license-function` repository.
   - Pull the latest changes:
     ```bash
     git pull origin main
     ```

2. **Navigate to the License Function Directory:**
   - Move to the `grsf-license-function` directory:
     ```bash
     cd grsf-license-function
     ```

3. **Build the Docker Image:**
   - Open a terminal and run the following command to build the Docker image:
     ```bash
     docker buildx build --platform "linux/amd64,linux/arm64" . --tag=runtime
     ```
     Note: the process is not working for me, when building multi-platform image on my arm64 based machine.
           I have to build the image for amd64 and then download it as a tar and use the tar to build the package
     ```bash
     docker buildx build --platform "linux/amd64" . --tag=runtime
     docker save runtime -o runtime.tar
     ```

4. **Build the Crossplane Package:**
   - Open a terminal and run the following command to build the crossplane function:
     ```bash
     crossplane xpkg build -f package --embed-runtime-image=image
     ```
     Note: the process is not working for me, when building multi-platform image on my arm64 based machine.
           I have to build the image for amd64 and then download it as a tar and use the tar to build the package
     ```bash
     docker buildx build --platform "linux/amd64" . --tag=runtime
     docker save runtime -o runtime.tar
     crossplane xpkg build -f package --embed-runtime-image-tarball=runtime.tar
     ```

5. **Publish the Crossplane Package:**
   - Open a terminal and run the following command to publish the crossplane function:
     ```bash
     crossplane xpkg push -f package/*.xpkg docker.io/grpl/grsf-lic:v0.1.6
     ```


