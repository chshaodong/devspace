package targetselector

import (
	"strings"

	"github.com/devspace-cloud/devspace/pkg/devspace/config/configutil"
	"github.com/devspace-cloud/devspace/pkg/devspace/config/versions/latest"
	"github.com/devspace-cloud/devspace/pkg/devspace/kubectl"
)

// SelectorParameter holds the information from the config and the command overrides
type SelectorParameter struct {
	ConfigParameter ConfigParameter
	CmdParameter    CmdParameter
}

// CmdParameter holds the parameter we receive from the command
type CmdParameter struct {
	Selector      string
	LabelSelector string
	Namespace     string
	ContainerName string
	PodName       string
	Pick          *bool
}

// ConfigParameter holds the parameter we receive from the config
type ConfigParameter struct {
	Selector      string
	LabelSelector map[string]string
	Namespace     string
	ContainerName string
}

// GetNamespace retrieves the target namespace
func (t *SelectorParameter) GetNamespace(config *latest.Config, kubeClient *kubectl.Client) (string, error) {
	if t.CmdParameter.Namespace != "" {
		return t.CmdParameter.Namespace, nil
	}
	if t.ConfigParameter.Namespace != "" {
		return t.ConfigParameter.Namespace, nil
	}
	if t.ConfigParameter.Selector != "" {
		selector, err := configutil.GetSelector(config, t.ConfigParameter.Selector)
		if err != nil {
			return "", err
		}
		if selector.Namespace != "" {
			return selector.Namespace, nil
		}
	}

	return kubeClient.Namespace, nil
}

// GetLabelSelector retrieves the label selector of the target
func (t *SelectorParameter) GetLabelSelector(config *latest.Config) (string, error) {
	if t.CmdParameter.LabelSelector != "" {
		return t.CmdParameter.LabelSelector, nil
	}
	if t.ConfigParameter.LabelSelector != nil {
		labelSelector := labelSelectorMapToString(t.ConfigParameter.LabelSelector)
		return labelSelector, nil
	}
	if t.ConfigParameter.Selector != "" {
		selector, err := configutil.GetSelector(config, t.ConfigParameter.Selector)
		if err != nil {
			return "", err
		}
		if selector.LabelSelector != nil {
			labelSelector := labelSelectorMapToString(selector.LabelSelector)
			return labelSelector, nil
		}
	}

	// We get the first selector if it exists
	if config != nil {
		if config.Dev != nil && config.Dev.Selectors != nil {
			selectors := config.Dev.Selectors
			if len(selectors) == 1 && selectors[0].LabelSelector != nil {
				labelSelector := labelSelectorMapToString(selectors[0].LabelSelector)
				return labelSelector, nil
			}
		}
	}

	return "", nil
}

func labelSelectorMapToString(m map[string]string) string {
	labels := make([]string, 0, len(m)-1)
	for key, value := range m {
		labels = append(labels, key+"="+value)
	}

	return strings.Join(labels, ",")
}

// GetPodName retrieves the pod name from the parameters
func (t *SelectorParameter) GetPodName() string {
	if t.CmdParameter.PodName != "" {
		return t.CmdParameter.PodName
	}

	return ""
}

// GetContainerName retrieves the container name from the parameters
func (t *SelectorParameter) GetContainerName() string {
	if t.CmdParameter.ContainerName != "" {
		return t.CmdParameter.ContainerName
	}
	if t.ConfigParameter.ContainerName != "" {
		return t.ConfigParameter.ContainerName
	}

	return ""
}
