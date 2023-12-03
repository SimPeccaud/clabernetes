package config

import (
	"fmt"
	"os"
	"strings"

	clabernetesconstants "github.com/srl-labs/clabernetes/constants"

	clabernetesapisv1alpha1 "github.com/srl-labs/clabernetes/apis/v1alpha1"
	claberneteserrors "github.com/srl-labs/clabernetes/errors"
	"gopkg.in/yaml.v3"
	k8scorev1 "k8s.io/api/core/v1"
	sigsyaml "sigs.k8s.io/yaml"
)

type bootstrapConfig struct {
	mergeMode                   string
	globalAnnotations           map[string]string
	globalLabels                map[string]string
	resourcesDefault            *k8scorev1.ResourceRequirements
	resourcesByContainerlabKind map[string]map[string]*k8scorev1.ResourceRequirements
	privilegedLauncher          bool
	containerlabDebug           bool
	inClusterDNSSuffix          string
	imagePullThroughMode        string
	launcherImage               string
	launcherImagePullPolicy     string
	launcherLogLevel            string
}

func bootstrapFromConfigMap( //nolint:gocyclo,funlen
	inMap map[string]string,
) (*bootstrapConfig, error) {
	bc := &bootstrapConfig{
		mergeMode:               "merge",
		inClusterDNSSuffix:      clabernetesconstants.KubernetesDefaultInClusterDNSSuffix,
		imagePullThroughMode:    clabernetesconstants.ImagePullThroughModeAuto,
		launcherImage:           os.Getenv(clabernetesconstants.LauncherImageEnv),
		launcherImagePullPolicy: clabernetesconstants.KubernetesImagePullIfNotPresent,
		launcherLogLevel:        clabernetesconstants.Info,
	}

	var outErrors []string

	mergeMode, mergeModeOk := inMap["mergeMode"]
	if mergeModeOk {
		bc.mergeMode = mergeMode
	}

	globalAnnotationsData, globalAnnotationsOk := inMap["globalAnnotations"]
	if globalAnnotationsOk {
		err := yaml.Unmarshal([]byte(globalAnnotationsData), &bc.globalAnnotations)
		if err != nil {
			outErrors = append(outErrors, err.Error())
		}
	}

	globalLabelsData, globalLabelsOk := inMap["globalLabels"]
	if globalLabelsOk {
		err := yaml.Unmarshal([]byte(globalLabelsData), &bc.globalLabels)
		if err != nil {
			outErrors = append(outErrors, err.Error())
		}
	}

	resourcesDefaultData, resourcesDefaultOk := inMap["resourcesDefault"]
	if resourcesDefaultOk {
		err := sigsyaml.Unmarshal([]byte(resourcesDefaultData), &bc.resourcesDefault)
		if err != nil {
			outErrors = append(outErrors, err.Error())
		}
	}

	resourcesByKindData, resourcesByKindOk := inMap["resourcesByContainerlabKind"]
	if resourcesByKindOk {
		err := sigsyaml.Unmarshal([]byte(resourcesByKindData), &bc.resourcesDefault)
		if err != nil {
			outErrors = append(outErrors, err.Error())
		}
	}

	inPrivilegedLauncher, inPrivilegedLauncherOk := inMap["privilegedLauncher"]
	if inPrivilegedLauncherOk {
		if strings.EqualFold(inPrivilegedLauncher, clabernetesconstants.True) {
			bc.privilegedLauncher = true
		}
	}

	inContainerlabDebug, inContainerlabDebugOk := inMap["containerlabDebug"]
	if inContainerlabDebugOk {
		if strings.EqualFold(inContainerlabDebug, clabernetesconstants.True) {
			bc.containerlabDebug = true
		}
	}

	inClusterDNSSuffix, inClusterDNSSuffixOk := inMap["inClusterDNSSuffix"]
	if inClusterDNSSuffixOk {
		bc.inClusterDNSSuffix = inClusterDNSSuffix
	}

	imagePullThroughMode, imagePullThroughModeOk := inMap["imagePullThroughMode"]
	if imagePullThroughModeOk {
		bc.imagePullThroughMode = imagePullThroughMode
	}

	launcherImage, launcherImageOk := inMap["launcherImage"]
	if launcherImageOk {
		bc.launcherImage = launcherImage
	}

	launcherImagePullPolicy, launcherImagePullPolicyOk := inMap["launcherImagePullPolicy"]
	if launcherImagePullPolicyOk {
		bc.launcherImagePullPolicy = launcherImagePullPolicy
	}

	launcherLogLevel, launcherLogLevelOk := inMap["launcherLogLevel"]
	if launcherLogLevelOk {
		bc.launcherLogLevel = launcherLogLevel
	}

	var err error

	if len(outErrors) > 0 {
		errors := ""

		for idx, outError := range outErrors {
			errors += fmt.Sprintf("error %d '%s'", idx, outError)
		}

		err = fmt.Errorf("%w: %s", claberneteserrors.ErrParse, errors)
	}

	return bc, err
}

// MergeFromBootstrapConfig accepts a bootstrap config configmap and the instance of the global
// config CR and merges the bootstrap config data onto the CR. The merge operation is based on the
// config merge mode set in both the bootstrap config and the CR (with the CR setting taking
// precedence).
func MergeFromBootstrapConfig(
	bootstrapConfigMap *k8scorev1.ConfigMap,
	config *clabernetesapisv1alpha1.Config,
) error {
	bootstrap, err := bootstrapFromConfigMap(bootstrapConfigMap.Data)
	if err != nil {
		return err
	}

	if bootstrap.mergeMode == "overwrite" {
		mergeFromBootstrapConfigReplace(bootstrap, config)
	} else {
		// should only ever be "merge" if it isn't "overwrite", but either way, fallback to merge...
		mergeFromBootstrapConfigMerge(bootstrap, config)
	}

	return nil
}

func mergeFromBootstrapConfigMerge(
	bootstrap *bootstrapConfig,
	config *clabernetesapisv1alpha1.Config,
) {
	for k, v := range bootstrap.globalAnnotations {
		_, exists := config.Spec.Metadata.Annotations[k]
		if exists {
			continue
		}

		config.Spec.Metadata.Annotations[k] = v
	}

	for k, v := range bootstrap.globalLabels {
		_, exists := config.Spec.Metadata.Labels[k]
		if exists {
			continue
		}

		config.Spec.Metadata.Labels[k] = v
	}

	if config.Spec.InClusterDNSSuffix == "" {
		config.Spec.InClusterDNSSuffix = bootstrap.inClusterDNSSuffix
	}

	if config.Spec.ImagePull.PullThroughOverride == "" {
		config.Spec.ImagePull.PullThroughOverride = bootstrap.imagePullThroughMode
	}

	if config.Spec.Deployment.ResourcesDefault == nil {
		config.Spec.Deployment.ResourcesDefault = bootstrap.resourcesDefault
	}

	for k, v := range bootstrap.resourcesByContainerlabKind {
		_, exists := config.Spec.Deployment.ResourcesByContainerlabKind[k]
		if exists {
			continue
		}

		config.Spec.Deployment.ResourcesByContainerlabKind[k] = v
	}

	if config.Spec.Deployment.LauncherImage == "" {
		config.Spec.Deployment.LauncherImage = bootstrap.launcherImage
	}

	if config.Spec.Deployment.LauncherImagePullPolicy == "" {
		config.Spec.Deployment.LauncherImagePullPolicy = bootstrap.launcherImagePullPolicy
	}

	if config.Spec.Deployment.LauncherLogLevel == "" {
		config.Spec.Deployment.LauncherLogLevel = bootstrap.launcherLogLevel
	}
}

func mergeFromBootstrapConfigReplace(
	bootstrap *bootstrapConfig,
	config *clabernetesapisv1alpha1.Config,
) {
	config.Spec = clabernetesapisv1alpha1.ConfigSpec{
		Metadata: clabernetesapisv1alpha1.ConfigMetadata{
			Annotations: bootstrap.globalAnnotations,
			Labels:      bootstrap.globalLabels,
		},
		InClusterDNSSuffix: bootstrap.inClusterDNSSuffix,
		ImagePull: clabernetesapisv1alpha1.ConfigImagePull{
			PullThroughOverride: bootstrap.imagePullThroughMode,
		},
		Deployment: clabernetesapisv1alpha1.ConfigDeployment{
			ResourcesDefault:            bootstrap.resourcesDefault,
			ResourcesByContainerlabKind: bootstrap.resourcesByContainerlabKind,
			PrivilegedLauncher:          bootstrap.privilegedLauncher,
			ContainerlabDebug:           bootstrap.containerlabDebug,
			LauncherImage:               bootstrap.launcherImage,
			LauncherImagePullPolicy:     bootstrap.launcherImagePullPolicy,
			LauncherLogLevel:            bootstrap.launcherLogLevel,
		},
	}
}
