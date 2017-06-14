/*-
 * Copyright (c) 2016,2017, F5 Networks, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package f5router

import (
	"os"
	"sort"
	"sync"

	"github.com/F5Networks/cf-bigip-ctlr/config"
	"github.com/F5Networks/cf-bigip-ctlr/registry"
	"github.com/F5Networks/cf-bigip-ctlr/registry/container"
	"github.com/F5Networks/cf-bigip-ctlr/route"
	"github.com/F5Networks/cf-bigip-ctlr/test_util"

	"code.cloudfoundry.org/routing-api/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("F5Router", func() {
	Describe("sorting", func() {
		Context("route config", func() {
			It("should sort correctly", func() {
				routeconfigs := routeConfigs{}

				expectedList := make(routeConfigs, 10)

				rc := routeConfig{}
				rc.Item.Backend.ServiceName = "bar"
				rc.Item.Backend.ServicePort = 80
				routeconfigs = append(routeconfigs, &rc)
				expectedList[1] = &rc

				rc = routeConfig{}
				rc.Item.Backend.ServiceName = "foo"
				rc.Item.Backend.ServicePort = 2
				routeconfigs = append(routeconfigs, &rc)
				expectedList[5] = &rc

				rc = routeConfig{}
				rc.Item.Backend.ServiceName = "foo"
				rc.Item.Backend.ServicePort = 8080
				routeconfigs = append(routeconfigs, &rc)
				expectedList[7] = &rc

				rc = routeConfig{}
				rc.Item.Backend.ServiceName = "baz"
				rc.Item.Backend.ServicePort = 1
				routeconfigs = append(routeconfigs, &rc)
				expectedList[2] = &rc

				rc = routeConfig{}
				rc.Item.Backend.ServiceName = "foo"
				rc.Item.Backend.ServicePort = 80
				routeconfigs = append(routeconfigs, &rc)
				expectedList[6] = &rc

				rc = routeConfig{}
				rc.Item.Backend.ServiceName = "foo"
				rc.Item.Backend.ServicePort = 9090
				routeconfigs = append(routeconfigs, &rc)
				expectedList[9] = &rc

				rc = routeConfig{}
				rc.Item.Backend.ServiceName = "baz"
				rc.Item.Backend.ServicePort = 1000
				routeconfigs = append(routeconfigs, &rc)
				expectedList[3] = &rc

				rc = routeConfig{}
				rc.Item.Backend.ServiceName = "foo"
				rc.Item.Backend.ServicePort = 8080
				routeconfigs = append(routeconfigs, &rc)
				expectedList[8] = &rc

				rc = routeConfig{}
				rc.Item.Backend.ServiceName = "foo"
				rc.Item.Backend.ServicePort = 1
				routeconfigs = append(routeconfigs, &rc)
				expectedList[4] = &rc

				rc = routeConfig{}
				rc.Item.Backend.ServiceName = "bar"
				rc.Item.Backend.ServicePort = 1
				routeconfigs = append(routeconfigs, &rc)
				expectedList[0] = &rc

				sort.Sort(routeconfigs)

				for i := range expectedList {
					Expect(routeconfigs[i]).To(Equal(expectedList[i]),
						"Sorted list elements should be equal")
				}
			})
		})

		Context("rules", func() {
			It("should sort correctly", func() {
				l7 := rules{}

				expectedList := make(rules, 10)

				p := rule{}
				p.FullURI = "bar"
				l7 = append(l7, &p)
				expectedList[1] = &p

				p = rule{}
				p.FullURI = "foo"
				l7 = append(l7, &p)
				expectedList[5] = &p

				p = rule{}
				p.FullURI = "foo"
				l7 = append(l7, &p)
				expectedList[7] = &p

				p = rule{}
				p.FullURI = "baz"
				l7 = append(l7, &p)
				expectedList[2] = &p

				p = rule{}
				p.FullURI = "foo"
				l7 = append(l7, &p)
				expectedList[6] = &p

				p = rule{}
				p.FullURI = "foo"
				l7 = append(l7, &p)
				expectedList[9] = &p

				p = rule{}
				p.FullURI = "baz"
				l7 = append(l7, &p)
				expectedList[3] = &p

				p = rule{}
				p.FullURI = "foo"
				l7 = append(l7, &p)
				expectedList[8] = &p

				p = rule{}
				p.FullURI = "foo"
				l7 = append(l7, &p)
				expectedList[4] = &p

				p = rule{}
				p.FullURI = "bar"
				l7 = append(l7, &p)
				expectedList[0] = &p

				sort.Sort(l7)

				for i := range expectedList {
					Expect(l7[i]).To(Equal(expectedList[i]),
						"Sorted list elements should be equal")
				}
			})
		})
	})

	Describe("verify configs", func() {
		It("should process the config correctly", func() {
			logger := test_util.NewTestZapLogger("router-test")
			mw := &MockWriter{}
			c := config.DefaultConfig()

			r, err := NewF5Router(logger, nil, mw)
			Expect(r).To(BeNil())
			Expect(err).To(HaveOccurred())

			r, err = NewF5Router(logger, c, nil)
			Expect(r).To(BeNil())
			Expect(err).To(HaveOccurred())

			c.BigIP.URL = "http://example.com"
			r, err = NewF5Router(logger, c, mw)
			Expect(r).To(BeNil())
			Expect(err).To(HaveOccurred())

			c.BigIP.User = "admin"
			r, err = NewF5Router(logger, c, mw)
			Expect(r).To(BeNil())
			Expect(err).To(HaveOccurred())

			c.BigIP.Pass = "pass"
			r, err = NewF5Router(logger, c, mw)
			Expect(r).To(BeNil())
			Expect(err).To(HaveOccurred())

			c.BigIP.Partitions = []string{"cf"}
			r, err = NewF5Router(logger, c, mw)
			Expect(r).To(BeNil())
			Expect(err).To(HaveOccurred())

			c.BigIP.ExternalAddr = "127.0.0.1"
			r, err = NewF5Router(logger, c, mw)
			Expect(r).NotTo(BeNil())
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("running router", func() {
		var mw *MockWriter
		var router *F5Router
		var err error
		var logger *test_util.TestZapLogger
		var c *config.Config

		BeforeEach(func() {
			logger = test_util.NewTestZapLogger("router-test")
			c = makeConfig()
			mw = &MockWriter{}

			router, err = NewF5Router(logger, c, mw)

			Expect(router).NotTo(BeNil(), "%v", router)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			if nil != logger {
				logger.Close()
			}
		})

		It("should run", func() {
			done := make(chan struct{})
			os := make(chan os.Signal)
			ready := make(chan struct{})

			go func() {
				defer GinkgoRecover()
				Expect(func() {
					err = router.Run(os, ready)
					Expect(err).NotTo(HaveOccurred())
					close(done)
				}).NotTo(Panic())
			}()
			// wait for the router to be ready
			Eventually(ready).Should(BeClosed(), "timed out waiting for ready")
			// send a kill signal to the router
			os <- MockSignal(123)
			//wait for the router to stop
			Eventually(done).Should(BeClosed(), "timed out waiting for Run to complete")
		})

		It("should update routes", func() {
			data := createTrie()
			done := make(chan struct{})
			os := make(chan os.Signal)
			ready := make(chan struct{})
			update := 0

			go func() {
				defer GinkgoRecover()
				Expect(func() {
					err = router.Run(os, ready)
					Expect(err).NotTo(HaveOccurred())
					close(done)
				}).NotTo(Panic())
			}()

			data.EachNodeWithPool(func(t *container.Trie) {
				t.Pool.Each(func(e *route.Endpoint) {
					go func(t *container.Trie, uri string) {
						router.RouteUpdate(
							registry.Add,
							t,
							route.Uri(uri),
						)
					}(data, t.ToPath())
				})
			})

			Eventually(mw.Input).Should(MatchJSON(expectedConfigs[update]))
			update++

			// make some changes and update the verification function
			p := data.Find(route.Uri("bar.cf.com"))
			Expect(p).NotTo(BeNil())
			removed := p.Remove(makeEndpoint("127.0.1.1"))
			Expect(removed).To(BeTrue())

			p = data.Find(route.Uri("baz.cf.com/segment1"))
			Expect(p).NotTo(BeNil())
			removed = p.Remove(makeEndpoint("127.0.3.2"))
			Expect(removed).To(BeTrue())

			p = data.Find(route.Uri("baz.cf.com"))
			Expect(p).NotTo(BeNil())
			added := p.Put(makeEndpoint("127.0.2.2"))
			Expect(added).To(BeTrue())

			removed = data.Delete(route.Uri("*.foo.cf.com"))
			Expect(removed).To(BeTrue())

			removed = data.Delete(route.Uri("foo.cf.com"))
			Expect(removed).To(BeTrue())

			router.RouteUpdate(
				registry.Remove,
				data,
				route.Uri("*.foo.cf.com"),
			)

			router.RouteUpdate(
				registry.Remove,
				data,
				route.Uri("foo.cf.com"),
			)

			Eventually(mw.Input).Should(MatchJSON(expectedConfigs[update]))
			update++

			p = route.NewPool(1, "qux.cf.com")
			p.Put(makeEndpoint("127.0.7.1"))
			data.Insert(route.Uri("qux.cf.com"), p)

			router.RouteUpdate(
				registry.Add,
				data,
				route.Uri("qux.cf.com"),
			)

			Eventually(mw.Input).Should(MatchJSON(expectedConfigs[update]))

			os <- MockSignal(123)
			Eventually(done).Should(BeClosed(), "timed out waiting for Run to complete")
		})

		It("should handle ssl and health monitors", func() {
			data := createTrie()
			done := make(chan struct{})
			os := make(chan os.Signal)
			ready := make(chan struct{})

			c.BigIP.HealthMonitors = []string{"Common/potato"}
			c.BigIP.SSLProfiles = []string{"Common/clientssl"}
			c.BigIP.Profiles = []string{"Common/http", "/Common/fakeprofile"}

			router, err = NewF5Router(logger, c, mw)

			go func() {
				defer GinkgoRecover()
				Expect(func() {
					err = router.Run(os, ready)
					Expect(err).NotTo(HaveOccurred())
					close(done)
				}).NotTo(Panic())
			}()
			Eventually(ready).Should(BeClosed())
			data.EachNodeWithPool(func(t *container.Trie) {
				t.Pool.Each(func(e *route.Endpoint) {
					go func(t *container.Trie, uri string) {
						router.RouteUpdate(
							registry.Add,
							t,
							route.Uri(uri),
						)
					}(data, t.ToPath())
				})
			})

			Eventually(mw.Input).Should(MatchJSON(expectedConfigs[3]))
		})

		Context("fail cases", func() {
			It("should error when not passing a URI for route update", func() {
				data := createTrie()
				done := make(chan struct{})
				os := make(chan os.Signal)
				ready := make(chan struct{})

				go func() {
					defer GinkgoRecover()
					Expect(func() {
						err = router.Run(os, ready)
						Expect(err).NotTo(HaveOccurred())
						close(done)
					}).NotTo(Panic())
				}()

				router.RouteUpdate(
					registry.Remove,
					data,
					route.Uri(""),
				)

				Eventually(logger).Should(Say("f5router-skipping-update"))
			})
		})

	})
})

type testRoutes struct {
	Key         route.Uri
	Addrs       []*route.Endpoint
	ContextPath string
}

type MockWriter struct {
	sync.Mutex
	input []byte
}

func (mw *MockWriter) GetOutputFilename() string {
	return "mock-file"
}

func (mw *MockWriter) Write(input []byte) (n int, err error) {
	mw.Lock()
	defer mw.Unlock()
	mw.input = input

	return len(input), nil
}

func (mw *MockWriter) Input() []byte {
	mw.Lock()
	defer mw.Unlock()
	dest := make([]byte, len(mw.input))
	l := copy(dest, mw.input)
	Expect(len(mw.input)).To(Equal(l))
	return dest
}

type MockSignal int

func (ms MockSignal) String() string {
	return "mock signal"
}

func (ms MockSignal) Signal() {
	return
}

func makeConfig() *config.Config {
	c := config.DefaultConfig()
	c.BigIP.URL = "http://example.com"
	c.BigIP.User = "admin"
	c.BigIP.Pass = "pass"
	c.BigIP.Partitions = []string{"cf"}
	c.BigIP.ExternalAddr = "127.0.0.1"

	return c
}

func makeEndpoints(addrs ...string) []*route.Endpoint {
	var r []*route.Endpoint
	for _, addr := range addrs {
		r = append(r, makeEndpoint(addr))
	}

	return r
}

func makeEndpoint(addr string) *route.Endpoint {
	r := route.NewEndpoint("1",
		addr,
		80,
		"1",
		"1",
		make(map[string]string),
		1,
		"",
		models.ModificationTag{
			Guid:  "1",
			Index: 1,
		},
	)
	return r
}

func createTrie() *container.Trie {
	data := container.NewTrie()

	routes := []testRoutes{
		{
			Key:         "foo.cf.com",
			Addrs:       makeEndpoints("127.0.0.1"),
			ContextPath: "/",
		},
		{
			Key:         "bar.cf.com",
			Addrs:       makeEndpoints("127.0.1.1", "127.0.1.2"),
			ContextPath: "/",
		},
		{
			Key:         "baz.cf.com",
			Addrs:       makeEndpoints("127.0.2.1"),
			ContextPath: "/",
		},
		{
			Key:         "baz.cf.com/segment1",
			Addrs:       makeEndpoints("127.0.3.1", "127.0.3.2"),
			ContextPath: "/segment1",
		},
		{
			Key:         "baz.cf.com/segment1/segment2/segment3",
			Addrs:       makeEndpoints("127.0.4.1", "127.0.4.2"),
			ContextPath: "/segment1/segment2/segment3",
		},
		{
			Key:         "*.cf.com",
			Addrs:       makeEndpoints("127.0.5.1"),
			ContextPath: "/",
		},
		{
			Key:         "*.foo.cf.com",
			Addrs:       makeEndpoints("127.0.6.1"),
			ContextPath: "/",
		},
	}

	for i := range routes {
		pool := route.NewPool(1, routes[i].ContextPath)
		for j := range routes[i].Addrs {
			pool.Put(routes[i].Addrs[j])
		}
		data.Insert(routes[i].Key, pool)
	}

	return data
}