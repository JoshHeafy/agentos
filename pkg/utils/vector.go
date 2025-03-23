package utils

import (
	"fmt"
	"log"
	"math"
	"strings"
)

func FloatSliceToPgvector(slice []float64) string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, val := range slice {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("%f", val))
	}
	sb.WriteString("]")
	return sb.String()
}

func ParsePgvector(vectorStr string) []float64 {
	var vec []float64
	vectorStr = strings.Trim(vectorStr, "[]")
	elements := strings.Split(vectorStr, ",")
	for _, elem := range elements {
		var value float64
		fmt.Sscanf(elem, "%f", &value)
		vec = append(vec, value)
	}
	return vec
}

func CosineSimilarity(vec1, vec2 []float64) float64 {
	if len(vec1) != len(vec2) {
		log.Fatal("Los vectores deben tener la misma longitud")
	}
	var dotProduct, normA, normB float64
	for i := range vec1 {
		dotProduct += vec1[i] * vec2[i]
		normA += vec1[i] * vec1[i]
		normB += vec2[i] * vec2[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
