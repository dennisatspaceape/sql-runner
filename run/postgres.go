//
// Copyright (c) 2015 Snowplow Analytics Ltd. All rights reserved.
//
// This program is licensed to you under the Apache License Version 2.0,
// and you may not use this file except in compliance with the Apache License Version 2.0.
// You may obtain a copy of the Apache License Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0.
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the Apache License Version 2.0 is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the Apache License Version 2.0 for the specific language governing permissions and limitations there under.
//
package run

import (
	"github.com/snowplow/sql-runner/playbook"
	"gopkg.in/pg.v2"
	"time"
)

// For Redshift queries
const (
	dialTimeout = 10 * time.Second
	readTimeout = 8 * time.Hour // TODO: make this user configurable
)

type PostgresTarget struct {
	playbook.Target
	Client *pg.DB
}

func NewPostgresTarget(target playbook.Target) *PostgresTarget {
	db := pg.Connect(&pg.Options{
		Host:        target.Host,
		Port:        target.Port,
		User:        target.Username,
		Password:    target.Password,
		Database:    target.Database,
		SSL:         false, // SSL is a problem for Redshift
		DialTimeout: dialTimeout,
		ReadTimeout: readTimeout,
	})

	return &PostgresTarget{target, db}
}

func (pt PostgresTarget) GetTarget() playbook.Target {
	return pt.Target
}

// Run a query against the target
func (pt PostgresTarget) RunQuery(query playbook.Query, queryPath string, variables map[string]interface{}) QueryStatus {

	script, err := prepareQuery(queryPath, query.Template, variables)
	if err != nil {
		return QueryStatus{query, queryPath, 0, err}
	}

	res, err := pt.Client.Exec(script)
	affected := 0
	if err == nil {
		affected = res.Affected()
	}

	return QueryStatus{query, queryPath, affected, err}
}
