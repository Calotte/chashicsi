# chashicsi
chashicsi is a sample csi used to help you create your own csi as a startup.

When pvc was created the CreateVolume in controller.go will be triggered

When pod that reference this pvc was created the NodePublishVolume in node.go will be invoked

And when pod was deleted the NodeUnpublishVolume was called as the log show
