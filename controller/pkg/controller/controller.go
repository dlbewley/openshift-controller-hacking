package controller

import (
        "fmt"

        osclient "github.com/openshift/origin/pkg/client"

        kclient "k8s.io/kubernetes/pkg/client/unversioned"
        kapi "k8s.io/kubernetes/pkg/api"
)

// Define an object for our controller to hold references to
// our OpenShift client
type Controller struct {
        openshiftClient *osclient.Client
        kubeClient *kclient.Client
}

// Function to instantiate a controller
func NewController(os *osclient.Client, ) *Controller {
        return &Controller{
                openshiftClient: os,
                kubeClient:      kc,
        }
}

// Our main function call
func (c *Controller) Run() {
        // Get a list of all the projects (namespaces) in the cluster
        // using the OpenShift client
        projects, err := c.openshiftClient.Projects().List(kapi.ListOptions{})
        if err != nil {
                fmt.Println(err)
        }

        // Iterate through the list of projects
        for _, project := range projects.Items {
                fmt.Printf("%s\n", project.ObjectMeta.Name)
        }
}
