package main

import (
	"encoding/json"
	"flag"
	"os"
	"strings"

	"github.com/mayflower/docker-ls/cli/docker-ls/response"
	"github.com/mayflower/docker-ls/cli/util"
	"github.com/mayflower/docker-ls/lib"
)

type tagDetailsCmd struct {
	flags *flag.FlagSet
}

func (r *tagDetailsCmd) execute(argv []string) (err error) {
	libCfg := lib.NewConfig()
	libCfg.BindToFlags(r.flags)

	cfg := newConfig()
	cfg.bindToFlags(r.flags,
		OPTION_PROGRESS|OPTION_JSON_OUTPUT|OPTION_MANIFEST_VERSION|OPTION_INTERACTIVE_PASSWORD)

	rawManifest := false
	r.flags.BoolVar(&rawManifest, "raw-manifest", rawManifest, "output raw manifest")

	parseHistory := false
	r.flags.BoolVar(&parseHistory, "parse-history", parseHistory, "try to parse history in raw manifest")

	err = r.flags.Parse(argv)
	if err != nil {
		return
	}

	if cfg.interactivePassword {
		err = util.PromptPassword(&libCfg)
		if err != nil {
			return
		}
	}

	args := r.flags.Args()
	if len(args) != 1 {
		r.flags.Usage()
		os.Exit(1)
	}

	ref := lib.EmptyRefspec()
	err = ref.Set(args[0])
	if err != nil {
		return
	}

	progress := NewProgressIndicator(cfg)
	progress.Start("requesting manifest")

	registryApi, err := lib.NewRegistryApi(libCfg)
	if err != nil {
		return
	}

	tagDetails, err := registryApi.GetTagDetails(ref, cfg.manifestVersion)

	progress.Progress()
	progress.Finish("done")

	if err != nil {
		return
	}

	if rawManifest {
		manifest := tagDetails.RawManifest()

		if parseHistory {
			r.parseHistory(manifest)
		}

		err = serializeToStdout(manifest, cfg)
	} else {
		err = serializeToStdout(response.NewTagDetailsL0(tagDetails, true), cfg)
	}

	return
}

func (r *tagDetailsCmd) parseHistory(rawManifest interface{}) {
	var manifest map[string]interface{}
	var ok bool

	if manifest, ok = rawManifest.(map[string]interface{}); !ok {
		return
	}

	var history []interface{}

	if _, ok = manifest["history"]; !ok {
		return
	}
	if history, ok = manifest["history"].([]interface{}); !ok {
		return
	}

	for _, rawItem := range history {
		var item map[string]interface{}

		if item, ok = rawItem.(map[string]interface{}); !ok {
			continue
		}
		if _, ok = item["v1Compatibility"]; !ok {
			continue
		}

		var rawHistoryItem string

		if rawHistoryItem, ok = item["v1Compatibility"].(string); !ok {
			continue
		}

		decoder := json.NewDecoder(strings.NewReader(rawHistoryItem))
		var decodedHistoryItem interface{}

		if decoder.Decode(&decodedHistoryItem) != nil {
			continue
		}

		item["v1Compatibility"] = decodedHistoryItem
	}
}

func newTagDetailsCmd(name string) (cmd *tagDetailsCmd) {
	cmd = &tagDetailsCmd{
		flags: flag.NewFlagSet(name, flag.ExitOnError),
	}

	cmd.flags.Usage = commandUsage(name, "<repository:reference>", "Inspect a singe tag.", cmd.flags)

	return
}
