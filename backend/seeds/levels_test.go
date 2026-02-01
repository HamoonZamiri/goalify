package seeds

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLevelsJSONValid(t *testing.T) {
	levels, err := loadLevels()
	require.NoError(t, err, "levels.json should parse successfully")
	require.NotEmpty(t, levels, "levels.json should contain at least one level")

	// Validate each level's structure and business rules
	for i, lvl := range levels {
		// IDs should be sequential starting from 1
		assert.Equal(t, int32(i+1), lvl.ID, "level IDs should be sequential")

		// XP and cash values must be positive
		assert.Greater(t, lvl.LevelUpXP, int32(0), "level %d: xp must be positive", lvl.ID)
		assert.Greater(t, lvl.CashReward, int32(0), "level %d: cash must be positive", lvl.ID)
	}
}

func TestLevelsJSONProgression(t *testing.T) {
	levels, err := loadLevels()
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(levels), 2, "need at least 2 levels to test progression")

	// Verify XP requirements generally increase (with some flexibility for milestones)
	for i := 1; i < len(levels); i++ {
		prev := levels[i-1]
		curr := levels[i]

		// XP should increase or stay the same (allowing for milestone adjustments)
		assert.GreaterOrEqual(t, curr.LevelUpXP, prev.LevelUpXP,
			"level %d xp (%d) should not be less than level %d xp (%d)",
			curr.ID, curr.LevelUpXP, prev.ID, prev.LevelUpXP)
	}
}

func TestLevelsCount(t *testing.T) {
	levels, err := loadLevels()
	require.NoError(t, err)

	// Verify we have exactly 100 levels as specified in the original migration
	assert.Equal(t, 100, len(levels), "should have exactly 100 levels")
}
