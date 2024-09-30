# chashicsi
chashicsi is a sample csi used to help you create your own csi as a startup.

When pvc was created the CreateVolume in controller.go will be triggered

When pod that reference this pvc was created the NodePublishVolume in node.go will be invoked

And when pod was deleted the NodeUnpublishVolume was called as the log show
![image](https://github.com/user-attachments/assets/a9851b32-9f5f-4398-8efa-bdc4c3be8a6c)

