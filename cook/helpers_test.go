package cook

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/gogrlx/grlx/types"
)

func TestExtractRequisites(t *testing.T) {
	// TODO: implement
}

func TestExtractIncludes(t *testing.T) {
	testCases := []struct {
		id          string
		sprout      string
		basepath    string
		recipe      types.RecipeName
		mapContents []string
	}{
		{
			id:          "dev",
			sprout:      "testSprout",
			basepath:    getBasePath(),
			recipe:      "dev",
			mapContents: []string{"apache", "missing"},
		},
		{
			id:          "independent",
			sprout:      "testSprout",
			basepath:    getBasePath(),
			recipe:      "independent",
			mapContents: []string{},
		},
		{
			id:          "apache init",
			sprout:      "testSprout",
			basepath:    getBasePath(),
			recipe:      "apache",
			mapContents: []string{"apache"},
		},
		{
			id:          "apache slash init",
			sprout:      "testSprout",
			basepath:    getBasePath(),
			recipe:      "apache.init.grlx",
			mapContents: []string{"apache"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.id, func(t *testing.T) {
			fp, err := ResolveRecipeFilePath(getBasePath(), tc.recipe)
			if err != nil {
				t.Error(err)
			}
			f, _ := os.ReadFile(fp)
			r, err := extractIncludes(tc.sprout, tc.basepath, string(tc.recipe), f)
			if err != nil {
				t.Error(err)
			}
			if len(r) != len(tc.mapContents) {
				t.Errorf("expected %v but got %v", tc.mapContents, r)
			}
			sort.Slice(r, func(i, j int) bool {
				return string(r[i]) < string(r[j])
			})
			sort.Strings(tc.mapContents)
			for i := range tc.mapContents {
				if string(r[i]) != tc.mapContents[i] {
					t.Errorf("expected %v but got %v", tc.mapContents, r)
				}
			}
		})
	}
}

func TestCollectAllIncludes(t *testing.T) {
	testCases := []struct {
		id     string
		recipe types.RecipeName
		sprout string
	}{{
		id:     "dev",
		recipe: "dev",
		sprout: "testSprout",
	}}
	for _, tc := range testCases {
		t.Run(tc.id, func(t *testing.T) {
			recipes, err := collectAllIncludes(tc.sprout, getBasePath(), tc.recipe)
			fmt.Printf("%v, %v", recipes, err)
		})
	}
}

func TestRelativeRecipeToAbsolute(t *testing.T) {
	testCases := []struct {
		id              string
		recipe          types.RecipeName
		filepath        string
		err             error
		relatedFilepath string
	}{{
		id:              "file doesn't exist",
		recipe:          "",
		filepath:        "",
		err:             os.ErrNotExist,
		relatedFilepath: "",
	}, {
		id:              "valid missing recipe",
		recipe:          ".missing",
		filepath:        "missing",
		err:             nil,
		relatedFilepath: filepath.Join(getBasePath(), "dev.grlx"),
	}}
	for _, tc := range testCases {
		t.Run(tc.id, func(t *testing.T) {
			filepath, err := relativeRecipeToAbsolute(getBasePath(), tc.relatedFilepath, tc.recipe)
			if string(filepath) != tc.filepath {
				t.Errorf("expected %s but got %s", tc.filepath, filepath)
			}
			if !errors.Is(err, tc.err) {
				t.Errorf("expected error %v but got %v", tc.err, err)
			}
		})
	}
}