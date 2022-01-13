package v1_test

import (
	"fmt"
	"regexp"
	"time"

	v1 "github.com/epinio/epinio/internal/api/v1"
	"github.com/gorilla/websocket"
	"k8s.io/apiserver/pkg/util/wsstream"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AppExec Endpoint", func() {
	var (
		appName   string
		namespace string
	)

	//containerImageURL := "splatform/sample-app"

	BeforeEach(func() {
		// namespace = catalog.NewNamespaceName()
		// env.SetupAndTargetNamespace(namespace)
		// appName = catalog.NewAppName()
		// env.MakeContainerImageApp(appName, 1, containerImageURL)
	})
	AfterEach(func() {
		//env.DeleteNamespace(namespace)
	})

	Describe("GET /namespaces/:namespace/applications/:app/exec", func() {
		var wsConn *websocket.Conn

		BeforeEach(func() {
			appName = "sample"
			namespace = "workspace"
			wsURL := fmt.Sprintf("%s%s/%s", websocketURL, v1.Root, v1.Routes.Path("AppExec", namespace, appName))
			wsConn = env.MakeWebSocketConnection(wsURL, wsstream.ChannelWebSocketProtocol)
		})

		FIt("runs a command and gets the output back", func() {
			// Run command: echo "test" > /workspace/test && exit
			// Check stdout stream (it should send back the command we sent)
			// Check if the file exists on the application Pod with kubectl

			// TODO: Do we want to test the base64 subprotocol used by the cli?

			wsConn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			wsConn.SetReadDeadline(time.Now().Add(10 * time.Second))

			message := ""
			r, err := regexp.Compile(`.*\$`) // Matches the bash command prompt
			Expect(err).ToNot(HaveOccurred())
			// Read until we reach the prompt
			for !r.MatchString(message) {
				_, messageBytes, err := wsConn.ReadMessage()
				Expect(err).ToNot(HaveOccurred())
				message = string(messageBytes[1:])
			}

			// Run the command
			//command := []byte("/bin/bash -c 'echo \"test\" > /workspace/test-echo && exit'")
			//command := append([]byte("0"), []byte(base64.URLEncoding.EncodeToString([]byte("ls\n\r")))...)

			// Beware! When the "raw" protocol is used (wsstream.ChannelWebSocketProtocol)
			// the channel is defined by the first byte.
			// In the wsstream.Base64ChannelWebSocketProtocol case, the first byte
			// is considered to be the ascii code of the channel. E.g. byte 48 for "0"
			// https://github.com/kubernetes/kubernetes/blob/46c5edbc58b81046ce799875dc611beaaf0ffb44/staging/src/k8s.io/apiserver/pkg/util/wsstream/conn.go#L261-L264
			command := append([]byte{0}, []byte([]byte("ls\r"))...)
			err = wsConn.WriteMessage(websocket.TextMessage, command)
			Expect(err).ToNot(HaveOccurred())

			_, messageBytes, err := wsConn.ReadMessage()
			Expect(err).NotTo(HaveOccurred())
			fmt.Printf("string(messageBytes) = %+v\n", string(messageBytes))
			Expect(string(messageBytes)).To(Equal("dimitris"))

			// err := wsConn.Close()
			// // With regular `ws` we could expect to not see any errors. With `wss`
			// // however, with a tls layer in the mix, we can expect to see a `broken
			// // pipe` issued. That is not a thing to act on, and is ignored.
			// if err != nil && strings.Contains(err.Error(), "broken pipe") {
			// 	return logs
			// }
			// Expect(err).ToNot(HaveOccurred())
		})
	})
})
