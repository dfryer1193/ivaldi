package merge

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

// YAML safely merges the src yaml into the dst yaml.
// Existing keys in dst are preserved, missing keys are added from src.
func YAML(dst, src []byte) ([]byte, error) {
	var dstNode yaml.Node
	if err := yaml.Unmarshal(dst, &dstNode); err != nil {
		return nil, err
	}

	var srcNode yaml.Node
	if err := yaml.Unmarshal(src, &srcNode); err != nil {
		return nil, err
	}

	if dstNode.Kind == yaml.DocumentNode && srcNode.Kind == yaml.DocumentNode {
		if len(dstNode.Content) > 0 && len(srcNode.Content) > 0 {
			mergeNodes(dstNode.Content[0], srcNode.Content[0])
		} else if len(dstNode.Content) == 0 && len(srcNode.Content) > 0 {
			dstNode.Content = append(dstNode.Content, srcNode.Content[0])
		}
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(&dstNode); err != nil {
		return nil, err
	}
	encoder.Close()

	return buf.Bytes(), nil
}

func mergeNodes(dst, src *yaml.Node) {
	if dst.Kind != yaml.MappingNode || src.Kind != yaml.MappingNode {
		return
	}

	for i := 0; i < len(src.Content); i += 2 {
		srcKey := src.Content[i]
		srcVal := src.Content[i+1]

		found := false
		for j := 0; j < len(dst.Content); j += 2 {
			dstKey := dst.Content[j]
			dstVal := dst.Content[j+1]

			if dstKey.Value == srcKey.Value {
				found = true
				if dstVal.Kind == yaml.MappingNode && srcVal.Kind == yaml.MappingNode {
					mergeNodes(dstVal, srcVal)
				}
				break
			}
		}

		if !found {
			dst.Content = append(dst.Content, srcKey, srcVal)
		}
	}
}
