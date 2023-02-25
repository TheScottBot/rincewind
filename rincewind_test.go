package rincewind

import (
	"github.com/stretchr/testify/assert"
	//nolint:staticcheck
	"testing"
)

func TestTranslate(t *testing.T) {
	t.Run("first", func(t *testing.T) {
		r := New()
		result, err := r.Translate(TranslationRequest{})

		assert.Equal(t, TranslationResponse{}, result)
		assert.NoError(t, err)
	})
}
