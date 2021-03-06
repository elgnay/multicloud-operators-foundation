package e2e

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	clusterinfov1beta1 "github.com/open-cluster-management/multicloud-operators-foundation/pkg/apis/cluster/v1beta1"
	"github.com/open-cluster-management/multicloud-operators-foundation/test/e2e/util"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

var clusterInfoGVR = schema.GroupVersionResource{
	Group:    "internal.open-cluster-management.io",
	Version:  "v1beta1",
	Resource: "managedclusterinfos",
}

var _ = ginkgo.Describe("Testing ManagedClusterInfo", func() {
	ginkgo.Context("Get ManagedClusterInfo", func() {
		ginkgo.It("should get a ManagedClusterInfo successfully in cluster", func() {
			gomega.Eventually(func() bool {
				exists, err := util.HasResource(dynamicClient, clusterInfoGVR, managedClusterName, managedClusterName)
				if err != nil {
					return false
				}
				return exists
			}, eventuallyTimeout, eventuallyInterval).Should(gomega.BeTrue())
		})

		ginkgo.It("should have a valid condition", func() {
			gomega.Eventually(func() bool {
				ManagedClusterInfo, err := util.GetResource(dynamicClient, clusterInfoGVR, managedClusterName, managedClusterName)
				if err != nil {
					return false
				}
				// check the ManagedClusterInfo status
				return util.GetConditionTypeFromStatus(ManagedClusterInfo, clusterinfov1beta1.ManagedClusterInfoSynced)
			}, eventuallyTimeout, eventuallyInterval).Should(gomega.BeTrue())
		})
	})

	ginkgo.Context("Delete ManagedClusterInfo Automatically after ManagedCluster is deleted", func() {
		var testManagedClusterName string

		ginkgo.BeforeEach(func() {
			managedCluster, err := util.CreateManagedCluster(dynamicClient)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			testManagedClusterName = managedCluster.GetName()

			//ManagedClusterinfo should exist
			gomega.Eventually(func() bool {
				existing, err := util.HasResource(dynamicClient, clusterInfoGVR, testManagedClusterName, testManagedClusterName)
				if err != nil {
					return false
				}
				return existing
			}, eventuallyTimeout, eventuallyInterval).Should(gomega.BeTrue())

			//Delete the managedcluster
			err = util.DeleteClusterResource(dynamicClient, util.ManagedClusterGVR, testManagedClusterName)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})

		ginkgo.It("clusterinfo should be deleted automitically.", func() {
			gomega.Eventually(func() bool {
				existing, err := util.HasResource(dynamicClient, clusterInfoGVR, testManagedClusterName, testManagedClusterName)
				if err != nil {
					return false
				}
				return existing
			}, eventuallyTimeout, eventuallyInterval).ShouldNot(gomega.BeTrue())
		})
	})
})
