package main

import (
	"bufio"
	"net"
	"net/http"
	"net/url"

	"github.com/alibaba/pouch/apis/types"
	"github.com/alibaba/pouch/test/request"

	"github.com/go-check/check"
)

// const defines common image name
const (
	busyboxImage    = "registry.hub.docker.com/library/busybox:latest"
	helloworldImage = "registry.hub.docker.com/library/hello-world"
)

// VerifyCondition is used to check the condition value.
type VerifyCondition func() bool

// SkipIfFalse skips the suite, if any of the conditions is not satisfied.
func SkipIfFalse(c *check.C, conditions ...VerifyCondition) {
	for _, con := range conditions {
		if con() == false {
			c.Skip("Skip test as condition is not matched")
		}
	}
}

// CreateBusyboxContainerOk creates a busybox container and asserts success.
func CreateBusyboxContainerOk(c *check.C, cname string, cmd ...string) {
	// If not specified, CMD executed in container is "top".
	if len(cmd) == 0 {
		cmd = []string{"top"}
	}

	resp, err := CreateBusyboxContainer(c, cname, cmd...)
	c.Assert(err, check.IsNil)
	CheckRespStatus(c, resp, 201)
}

// CreateBusyboxContainer creates a basic container using busybox image.
func CreateBusyboxContainer(c *check.C, cname string, cmd ...string) (*http.Response, error) {
	q := url.Values{}
	q.Add("name", cname)

	obj := map[string]interface{}{
		"Image":      busyboxImage,
		"Cmd":        cmd,
		"HostConfig": map[string]interface{}{},
	}

	path := "/containers/create"
	query := request.WithQuery(q)
	body := request.WithJSONBody(obj)
	return request.Post(path, query, body)
}

// StartContainerOk starts the container and asserts success.
func StartContainerOk(c *check.C, cname string) {
	resp, err := StartContainer(c, cname)
	c.Assert(err, check.IsNil)

	CheckRespStatus(c, resp, 204)
}

// StartContainer starts the container.
func StartContainer(c *check.C, cname string) (*http.Response, error) {
	return request.Post("/containers/" + cname + "/start")
}

// DelContainerForceOk forcely deletes the container and asserts success.
func DelContainerForceOk(c *check.C, cname string) {
	resp, err := DelContainerForce(c, cname)
	c.Assert(err, check.IsNil)

	CheckRespStatus(c, resp, 204)
}

// DelContainerForce forcely deletes the container.
func DelContainerForce(c *check.C, cname string) (*http.Response, error) {
	q := url.Values{}
	q.Add("force", "true")
	return request.Delete("/containers/"+cname, request.WithQuery(q))
}

// StopContainerOk stops the container and asserts success..
func StopContainerOk(c *check.C, cname string) {
	resp, err := StopContainer(c, cname)
	c.Assert(err, check.IsNil)

	CheckRespStatus(c, resp, 204)
}

// StopContainer stops the container.
func StopContainer(c *check.C, cname string) (*http.Response, error) {
	return request.Post("/containers/" + cname + "/stop")
}

// PauseContainerOk pauses the container and asserts success..
func PauseContainerOk(c *check.C, cname string) {
	resp, err := PauseContainer(c, cname)
	c.Assert(err, check.IsNil)

	CheckRespStatus(c, resp, 204)
}

// PauseContainer pauses the container.
func PauseContainer(c *check.C, cname string) (*http.Response, error) {
	return request.Post("/containers/" + cname + "/pause")
}

// UnpauseContainerOk unpauses the container and asserts success..
func UnpauseContainerOk(c *check.C, cname string) {
	resp, err := UnpauseContainer(c, cname)
	c.Assert(err, check.IsNil)

	CheckRespStatus(c, resp, 204)
}

// UnpauseContainer unpauses the container.
func UnpauseContainer(c *check.C, cname string) (*http.Response, error) {
	return request.Post("/containers/" + cname + "/unpause")
}

// CheckRespStatus checks the http.Response.Status is equal to status.
func CheckRespStatus(c *check.C, resp *http.Response, status int) {
	if resp.StatusCode != status {
		got := types.Error{}
		_ = request.DecodeBody(&got, resp.Body)
		c.Assert(resp.StatusCode, check.Equals, status, check.Commentf("Error:%s", got.Message))
	}
}

// IsContainerCreated returns true is container's state is created.
func IsContainerCreated(c *check.C, cname string) (bool, error) {
	return isContainerStateEqual(c, cname, "created")
}

// IsContainerRunning returns true is container's state is running.
func IsContainerRunning(c *check.C, cname string) (bool, error) {
	return isContainerStateEqual(c, cname, "running")
}

func isContainerStateEqual(c *check.C, cname string, status string) (bool, error) {
	resp, err := request.Get("/containers/" + cname + "/json")
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 200)

	defer resp.Body.Close()
	got := types.ContainerJSON{}
	err = request.DecodeBody(&got, resp.Body)
	c.Assert(err, check.IsNil)

	if got.State == nil {
		return false, nil
	}

	return string(got.State.Status) == status, nil
}

// DelNetworkOk deletes the network and asserts success.
func DelNetworkOk(c *check.C, cname string) {
	resp, err := DelNetwork(c, cname)
	c.Assert(err, check.IsNil)

	CheckRespStatus(c, resp, 204)
}

// DelNetwork  deletes the network.
func DelNetwork(c *check.C, cname string) (*http.Response, error) {
	return request.Delete("/networks/" + cname)
}

// CreateExecEchoOk exec process's environment with "echo" CMD.
func CreateExecEchoOk(c *check.C, cname string) string {
	// NOTICE:
	// All files in the obj is needed, or start a new process may hang.
	obj := map[string]interface{}{
		"Cmd":          []string{"echo", "test"},
		"Detach":       true,
		"AttachStderr": true,
		"AttachStdout": true,
		"AttachStdin":  true,
		"Privileged":   false,
		"User":         "",
	}
	body := request.WithJSONBody(obj)

	resp, err := request.Post("/containers/"+cname+"/exec", body)
	c.Assert(err, check.IsNil)
	CheckRespStatus(c, resp, 201)

	var got types.ExecCreateResp
	request.DecodeBody(&got, resp.Body)
	return got.ID
}

// StartContainerExecOk starts executing a process in the container and asserts success.
func StartContainerExecOk(c *check.C, execid string, tty bool, detach bool) {
	resp, conn, _, err := StartContainerExec(c, execid, tty, detach)
	c.Assert(err, check.IsNil)

	// TODO: fix to use 200
	CheckRespStatus(c, resp, 101)
	defer conn.Close()
}

// StartContainerExec starts executing a process in the container.
func StartContainerExec(c *check.C, execid string, tty bool, detach bool) (*http.Response, net.Conn, *bufio.Reader, error) {

	obj := map[string]interface{}{
		"Detach": detach,
		"Tty":    tty,
	}
	body := request.WithJSONBody(obj)

	resp, conn, reader, err := request.Hijack("/exec/"+execid+"/start", body, request.WithHeader("Content-Type", "text/plain"))
	return resp, conn, reader, err
}
