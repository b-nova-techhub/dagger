package gl_betty_bonus_checker

import (
	"dagger.io/dagger"

	"dagger.io/dagger/core"
	"universe.dagger.io/go"
)

dagger.#Plan & {
	actions: {
		source: core.#Source & {
			path: "."
			exclude: [
				"build",
				"*.cue",
				"*.md",
				".git",
			]
		}

		build: {
			getCode: core.#Source & {
				path: "."
			}

			test: go.#Test & {
				source:  getCode.output
				package: "./..."
			}

			goBuild: go.#Build & {
				source: getCode.output
			}
		}
	}
}
