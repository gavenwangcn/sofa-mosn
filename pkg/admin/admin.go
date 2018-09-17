/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package admin

import (
	"encoding/json"
	"net/http"

	"github.com/alipay/sofa-mosn/pkg/log"
)

func configDump(w http.ResponseWriter, _ *http.Request) {
	if buf, err := json.Marshal(GetEffectiveConfig()); err == nil {
		w.Write(buf)
	} else {
		w.WriteHeader(500)
		w.Write([]byte(`{ error: "internal error" }`))
		log.DefaultLogger.Errorf("Admin API: ConfigDump failed, cause by %s", err)
	}
}

func Start(adminConfig AdminConfig, config interface{}) *http.Server {
	// merge MOSNConfig into global context
	var originalConf map[string]interface{}
	data, _ := json.Marshal(config)
	json.Unmarshal(data, &originalConf)
	Set("original_config", originalConf)

	addr := adminConfig.Address
	if addr == "" {
		addr = ":8888"
	}
	srv := &http.Server{Addr: addr}
	log.DefaultLogger.Infof("Admin server serve on %s", addr)

	go func() {
		http.HandleFunc("/api/v1/config_dump", configDump)
		if err := srv.ListenAndServe(); err != nil {
			log.DefaultLogger.Errorf("Admin server: ListenAndServe() error: %s", err)
		}
	}()

	return srv
}
