// Package mutate deals with AdmissionReview requests and responses, it takes in the request body and returns a readily converted JSON []byte that can be
// returned from a http Handler w/o needing to further convert or modify it, it also makes testing Mutate() kind of easy w/o need for a fake http server, etc.
package validate

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	v1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Validate validates
func Validate(body []byte, verbose bool) ([]byte, error) {
	if verbose {
		log.Printf("recv: %s\n", string(body))
	}

	// unmarshal request into AdmissionReview struct
	admReview := v1beta1.AdmissionReview{}
	if err := json.Unmarshal(body, &admReview); err != nil {
		return nil, fmt.Errorf("unmarshaling request failed with %s", err)
	}

	var err error
	var pod *corev1.Pod
	var validContainer bool

	responseBody := []byte{}
	ar := admReview.Request
	resp := v1beta1.AdmissionResponse{}
	if ar != nil {

		// get the Pod object and unmarshal it into its struct, if we cannot, we might as well stop here
		if err := json.Unmarshal(ar.Object.Raw, &pod); err != nil {
			return nil, fmt.Errorf("unable unmarshal pod json object %v", err)
		}
		// set response options
		resp.Allowed = true // allow pods by default
		resp.UID = ar.UID

		// add some audit annotations, helpful to know why a object was reviewed
		resp.AuditAnnotations = map[string]string{
			"reviewedResourceRequestsAndLimits": "true",
		}

		validatedContainers := []map[string]bool{}
		for i := range pod.Spec.Containers {
			// Only accept pods that are in the guaranteed QoS
			containerName := pod.Spec.Containers[i].Name
			containerCpuRequests := pod.Spec.Containers[i].Resources.Requests.Cpu().Value()
			containerCpuLimits := pod.Spec.Containers[i].Resources.Limits.Cpu().Value()
			containerMemoryRequests := pod.Spec.Containers[i].Resources.Requests.Memory().Value()
			containerMemoryLimits := pod.Spec.Containers[i].Resources.Limits.Memory().Value()
			log.Printf("Container %s, Requests: [CPU: %d, Memory: %d], Limits: [CPU: %d, Memory: %d]", containerName, containerCpuRequests, containerMemoryRequests, containerCpuLimits, containerMemoryLimits)
			if ((containerCpuRequests + containerCpuLimits + containerMemoryRequests + containerMemoryLimits) == 0 ) {
				log.Print("Container is in the BestEffort QoS. Marked as invalid.")
				validContainer = false
			} else if ((containerCpuRequests == containerCpuLimits) && (containerMemoryRequests == containerMemoryLimits)) {
				log.Print("Container is in the Guaranteed QoS. Marked as valid.")
				validContainer = true
			} else {
				log.Print("Container is in the Burstable QoS. Marked as invalid.")
				validContainer = false
			}
			container := map[string]bool{
				containerName: validContainer,
			}
			validatedContainers = append(validatedContainers, container)
			
		}

		// Get list of non-valid containers
		var nonValidContainers []string
		for i := range validatedContainers {
			for k, v := range validatedContainers[i] {
				log.Printf("Container: %s, valid: %t", k, v)	
				if (v == false) {
					nonValidContainers = append(nonValidContainers, k)
				}
			}
		}

		// If there is any non-valid container then reject the creation
		if (len(nonValidContainers) > 0 ) {
			resp.Allowed = false
			nonValidContainersNames := strings.Join(nonValidContainers, ", ")
			statusMessage := "The following non-valid containers prevented the pod creation: " + nonValidContainersNames
			log.Print(statusMessage)
			resp.Result = &metav1.Status{
				Message: statusMessage,
				Status: "Failure",
			}
		} else {
			// Success
			statusMessage := "The pod is valid, pod creation can proceed"
			log.Print(statusMessage)
			resp.Result = &metav1.Status{
				Message: statusMessage,
				Status: "Success",
			}
		}


		admReview.Response = &resp
		// back into JSON so we can return the finished AdmissionReview w/ Response directly
		// w/o needing to convert things in the http handler
		responseBody, err = json.Marshal(admReview)
		if err != nil {
			return nil, err // untested section
		}
	}

	if verbose {
		log.Printf("resp: %s\n", string(responseBody))
	}

	return responseBody, nil
}
