//go:build e2e
// +build e2e

package e2e

import (
	"os"
	"testing"

	eraserv1alpha1 "github.com/Azure/eraser/api/v1alpha1"
	"github.com/Azure/eraser/test/e2e/util"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
)

func TestMain(m *testing.M) {
	utilruntime.Must(eraserv1alpha1.AddToScheme(scheme.Scheme))

	removerImage := util.ParsedImages.RemoverImage
	managerImage := util.ParsedImages.ManagerImage
	collectorImage := util.ParsedImages.CollectorImage
	scannerImage := util.ParsedImages.ScannerImage

	util.Testenv = env.NewWithConfig(envconf.New())
	// Create KinD Cluster
	util.Testenv.Setup(
		envfuncs.CreateKindClusterWithConfig(util.KindClusterName, util.NodeVersion, util.KindConfigPath),
		envfuncs.CreateNamespace(util.TestNamespace),
		util.LoadImageToCluster(util.KindClusterName, util.ManagerImage, util.ManagerTarballPath),
		util.LoadImageToCluster(util.KindClusterName, util.RemoverImage, util.RemoverTarballPath),
		util.LoadImageToCluster(util.KindClusterName, util.CollectorImage, util.CollectorTarballPath),
		util.LoadImageToCluster(util.KindClusterName, util.ScannerImage, util.ScannerTarballPath),
		util.LoadImageToCluster(util.KindClusterName, util.RemoverImage, util.RemoverTarballPath),
		util.HelmDeployLatestEraserRelease(util.TestNamespace,
			"--set", util.ScannerEnable.Set("false"),
			"--set", util.CollectorEnable.Set("false"),
			"--set", util.RemoverImageRepo.Set(removerImage.Repo),
			"--set", util.RemoverImageTag.Set(removerImage.Tag),
			"--set", util.ManagerImageRepo.Set(managerImage.Repo),
			"--set", util.ManagerImageTag.Set(managerImage.Tag),
		),
		util.UpgradeEraserHelm(util.TestNamespace,
			"--set", util.ScannerEnable.Set("true"),
			"--set", util.ScannerImageRepo.Set(scannerImage.Repo),
			"--set", util.ScannerImageTag.Set(scannerImage.Tag),
			"--set", util.RemoverImageRepo.Set(removerImage.Repo),
			"--set", util.RemoverImageTag.Set(removerImage.Tag),
			"--set", util.CollectorEnable.Set("true"),
			"--set", util.CollectorImageRepo.Set(collectorImage.Repo),
			"--set", util.CollectorImageTag.Set(collectorImage.Tag),
			"--set", util.ManagerImageRepo.Set(managerImage.Repo),
			"--set", util.ManagerImageTag.Set(managerImage.Tag),
			"--set-json", util.ImagePullSecrets.Set(util.ImagePullSecretJSON),
			"--set", util.CleanupOnSuccessDelay.Set("1m"),
		),
	).Finish(
		envfuncs.DestroyKindCluster(util.KindClusterName),
	)
	os.Exit(util.Testenv.Run(m))
}
