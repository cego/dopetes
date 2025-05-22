package routines

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/cego/dopetes/model"
	"github.com/cego/go-lib"
	"github.com/elastic/go-elasticsearch/v9"
)

func PushDockerBuildxHistoryToElastic(logger cego.Logger, config *model.DopetesConfig, elasticClient *elasticsearch.Client, state *model.DockerBuildxHistoryState) error {
	out, err := exec.Command("docker", "buildx", "history", "ls", "--format=json", "--filter=status!=running").Output()
	if err != nil {
		return err
	}

	split := strings.Split(string(out), "\n")

	reMaterialUri := regexp.MustCompile(`pkg:docker/(.*)\?`)

	reHistoryLsId := regexp.MustCompile(".*/(.*)$")

	for _, line := range split {
		if line == "" {
			continue
		}
		dockerBuildxHistoryLs := &model.DockerBuildxHistoryLs{}
		err = json.Unmarshal([]byte(line), dockerBuildxHistoryLs)
		if err != nil {
			return err
		}

		matches := reHistoryLsId.FindStringSubmatch(dockerBuildxHistoryLs.Ref)
		id := matches[1]

		if state.HasId(id) {
			continue
		}

		state.AddId(id)

		out, err = exec.Command("docker", "buildx", "history", "inspect", id, "--format=json").Output()
		if err != nil {
			return err
		}

		dockerBuildxHistoryInspect := &model.DockerBuildxHistoryInspect{}
		err = json.Unmarshal(out, dockerBuildxHistoryInspect)
		if err != nil {
			return err
		}

		for _, material := range dockerBuildxHistoryInspect.Materials {
			matches = reMaterialUri.FindStringSubmatch(material.URI)
			imageRef := strings.ReplaceAll(matches[1], "@", ":")

			logger.Debug(fmt.Sprintf("Detected docker buildx history material %s", imageRef))

			document := &model.ElasticDocument{
				Timestamp: time.Now().Format(time.RFC3339),
				Message:   fmt.Sprintf("dopetes detected docker buildx history inspect material URI  %s", imageRef),
				ImageName: imageRef,
				EventRaw:  string(out),
			}
			data, _ := json.Marshal(document)
			_, err = elasticClient.Index(config.Elasticsearch.Index, bytes.NewReader(data))
			if err != nil {
				logger.Error(err.Error())
				continue
			}
		}

	}

	return nil
}

func StartDockerBuildxHistoryInterval(logger cego.Logger, config *model.DopetesConfig, elasticClient *elasticsearch.Client, state *model.DockerBuildxHistoryState) {
	go func() {
		for {
			err := PushDockerBuildxHistoryToElastic(logger, config, elasticClient, state)
			if err != nil {
				logger.Error(err.Error())
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
}
