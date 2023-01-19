package dagger

import (
	"dagger.io/dagger"

	"dagger.io/dagger/core"
	"universe.dagger.io/go"
)

dagger.#Plan & {
	client: filesystem: "./build": write: contents: actions.build.build.output
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
			name: "build"

			test: go.#Test & {
				source:  actions.source.output
				package: "./..."
			}

			build: go.#Build & {
				source: actions.source.output
			}
		}
	}
}
