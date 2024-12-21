package builder

type ContainerImage struct {
	ContainerName string
	Image         string
}

func (c *ContainerImage) String() string {
	return c.ContainerName + "," + c.Image
}

func ContainerImagesToStringArray(containerImages []ContainerImage) []string {
	results := make([]string, 0, len(containerImages))
	for _, containerImage := range containerImages {
		results = append(results, containerImage.String())
	}

	return results
}
