package project

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

const (
	inertiaDeployTest = "https://github.com/ubclaunchpad/inertia-deploy-test.git"
)

var urlVariations = []string{
	"git@github.com:ubclaunchpad/inertia.git",
	"https://github.com/ubclaunchpad/inertia.git",
	"git://github.com/ubclaunchpad/inertia.git",
	"git://github.com/ubclaunchpad/inertia",
}

func getInertiaDeployTestKey() (ssh.AuthMethod, error) {
	pemFile, err := os.Open("../../../test/keys/id_rsa")
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(pemFile)
	if err != nil {
		return nil, err
	}
	return ssh.NewPublicKeys("git", bytes, "")
}

func TestCloneIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	dir := "./test_clone/"
	repo, err := clone(dir, inertiaDeployTest, "dev", nil, os.Stdout)
	defer os.RemoveAll(dir)
	assert.Nil(t, err)

	head, err := repo.Head()
	assert.Nil(t, err)
	assert.Equal(t, "dev", head.Name().Short())
}

func TestForcePullIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	dir := "./test_force_pull/"
	auth, err := getInertiaDeployTestKey()
	assert.Nil(t, err)
	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: inertiaDeployTest,
	})
	defer os.RemoveAll(dir)
	assert.Nil(t, err)
	forcePulledRepo, err := forcePull(dir, repo, auth, os.Stdout)
	assert.Nil(t, err)

	// Try switching branches
	err = updateRepository(dir, forcePulledRepo, "dev", auth, os.Stdout)
	assert.Nil(t, err)
	err = updateRepository(dir, forcePulledRepo, "master", auth, os.Stdout)
	assert.Nil(t, err)
}

func TestUpdateRepositoryIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	dir := "./test_update/"
	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: inertiaDeployTest,
	})
	defer os.RemoveAll(dir)
	assert.Nil(t, err)

	// Try switching branches
	err = updateRepository(dir, repo, "master", nil, os.Stdout)
	assert.Nil(t, err)
	err = updateRepository(dir, repo, "dev", nil, os.Stdout)
	assert.Nil(t, err)
}