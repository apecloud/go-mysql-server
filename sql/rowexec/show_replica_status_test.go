// Copyright 2026 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rowexec

import (
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/binlogreplication"
	"github.com/dolthub/go-mysql-server/sql/plan"
)

type replicaStatusController struct {
	status *binlogreplication.ReplicaStatus
}

func (c *replicaStatusController) StartReplica(*sql.Context) error { return nil }
func (c *replicaStatusController) StopReplica(*sql.Context) error  { return nil }
func (c *replicaStatusController) SetReplicationSourceOptions(
	*sql.Context, []binlogreplication.ReplicationOption,
) error {
	return nil
}
func (c *replicaStatusController) SetReplicationFilterOptions(
	*sql.Context, []binlogreplication.ReplicationOption,
) error {
	return nil
}
func (c *replicaStatusController) GetReplicaStatus(
	*sql.Context,
) (*binlogreplication.ReplicaStatus, error) {
	return c.status, nil
}
func (c *replicaStatusController) ResetReplica(*sql.Context, bool) error { return nil }

func TestBuildShowReplicaStatusPreservesIgnoredSSLContract(t *testing.T) {
	lastIoErrorTimestamp := time.Date(2026, time.August, 24, 1, 2, 3, 0, time.UTC)
	lastSqlErrorTimestamp := time.Date(2026, time.August, 24, 4, 5, 6, 0, time.UTC)
	expectedColumnNames := []string{
		"Replica_IO_State", "Source_Host", "Source_User", "Source_Port", "Connect_Retry",
		"Source_Log_File", "Read_Source_Log_Pos", "Relay_Log_File", "Relay_Log_Pos", "Relay_Source_Log_File",
		"Replica_IO_Running", "Replica_SQL_Running", "Replicate_Do_DB", "Replicate_Ignore_DB", "Replicate_Do_Table",
		"Replicate_Ignore_Table", "Replicate_Wild_Do_Table", "Replicate_Wild_Ignore_Table", "Last_Errno", "Last_Error",
		"Skip_Counter", "Exec_Source_Log_Pos", "Relay_Log_Space", "Until_Condition", "Until_Log_File",
		"Until_Log_Pos", "Source_SSL_Allowed", "Source_SSL_CA_File", "Source_SSL_CA_Path", "Source_SSL_Cert",
		"Source_SSL_Cipher", "Source_SSL_CRL_File", "Source_SSL_CRL_Path", "Source_SSL_Key", "Source_SSL_Verify_Server_Cert",
		"Seconds_Behind_Source", "Last_IO_Errno", "Last_IO_Error", "Last_SQL_Errno", "Last_SQL_Error",
		"Replicate_Ignore_Server_Ids", "Source_Server_Id", "Source_UUID", "Source_Info_File", "SQL_Delay",
		"SQL_Remaining_Delay", "Replica_SQL_Running_State", "Source_Retry_Count", "Source_Bind", "Last_IO_Error_Timestamp",
		"Last_SQL_Error_Timestamp", "Retrieved_Gtid_Set", "Executed_Gtid_Set", "Auto_Position", "Replicate_Rewrite_DB",
	}
	expected := sql.Row{
		"",               // Replica_IO_State
		"source.example", // Source_Host
		"replicator",     // Source_User
		uint(3306),       // Source_Port
		uint32(7),        // Connect_Retry
		"INVALID",        // Source_Log_File
		0,                // Read_Source_Log_Pos
		nil,              // Relay_Log_File
		nil,              // Relay_Log_Pos
		"INVALID",        // Relay_Source_Log_File
		"Yes",            // Replica_IO_Running
		"Connecting",     // Replica_SQL_Running
		nil,              // Replicate_Do_DB
		nil,              // Replicate_Ignore_DB
		"db1.t1,db2.t2",  // Replicate_Do_Table
		"db3.t3",         // Replicate_Ignore_Table
		nil,              // Replicate_Wild_Do_Table
		nil,              // Replicate_Wild_Ignore_Table
		uint(22),         // Last_Errno
		"sql error",      // Last_Error
		nil,              // Skip_Counter
		0,                // Exec_Source_Log_Pos
		nil,              // Relay_Log_Space
		"None",           // Until_Condition
		nil,              // Until_Log_File
		nil,              // Until_Log_Pos
		"Ignored",        // Source_SSL_Allowed
		nil,              // Source_SSL_CA_File
		nil,              // Source_SSL_CA_Path
		nil,              // Source_SSL_Cert
		nil,              // Source_SSL_Cipher
		nil,              // Source_SSL_CRL_File
		nil,              // Source_SSL_CRL_Path
		nil,              // Source_SSL_Key
		nil,              // Source_SSL_Verify_Server_Cert
		0,                // Seconds_Behind_Source
		uint(11),         // Last_IO_Errno
		"io error",       // Last_IO_Error
		uint(22),         // Last_SQL_Errno
		"sql error",      // Last_SQL_Error
		nil,              // Replicate_Ignore_Server_Ids
		"42",             // Source_Server_Id
		"source-uuid",    // Source_UUID
		nil,              // Source_Info_File
		0,                // SQL_Delay
		0,                // SQL_Remaining_Delay
		nil,              // Replica_SQL_Running_State
		uint64(99),       // Source_Retry_Count
		nil,              // Source_Bind
		lastIoErrorTimestamp.Format(time.UnixDate),
		lastSqlErrorTimestamp.Format(time.UnixDate),
		"retrieved-gtid", // Retrieved_Gtid_Set
		"executed-gtid",  // Executed_Gtid_Set
		true,             // Auto_Position
		nil,              // Replicate_Rewrite_DB
	}

	ctx := sql.NewEmptyContext()
	node := plan.NewShowReplicaStatus()
	schema := node.Schema(ctx)
	actualColumnNames := make([]string, len(schema))
	sslAllowedIndex := -1
	for i, column := range schema {
		actualColumnNames[i] = column.Name
		if column.Name == "Source_SSL_Allowed" {
			sslAllowedIndex = i
		}
	}
	require.Equal(t, expectedColumnNames, actualColumnNames)
	require.Equal(t, 26, sslAllowedIndex)
	var rows []sql.Row
	for _, sourceSSL := range []bool{false, true} {
		t.Run(map[bool]string{false: "SourceSsl false", true: "SourceSsl true"}[sourceSSL], func(t *testing.T) {
			status := &binlogreplication.ReplicaStatus{
				LastSqlErrorTimestamp: &lastSqlErrorTimestamp,
				LastIoErrorTimestamp:  &lastIoErrorTimestamp,
				LastSqlError:          "sql error",
				LastIoError:           "io error",
				SourceHost:            "source.example",
				SourceUser:            "replicator",
				SourceServerId:        "42",
				SourceServerUuid:      "source-uuid",
				RetrievedGtidSet:      "retrieved-gtid",
				ExecutedGtidSet:       "executed-gtid",
				ReplicaIoRunning:      "Yes",
				ReplicaSqlRunning:     "Connecting",
				ReplicateDoTables:     []string{"db1.t1", "db2.t2"},
				ReplicateIgnoreTables: []string{"db3.t3"},
				SourceRetryCount:      99,
				SourcePort:            3306,
				LastIoErrNumber:       11,
				LastSqlErrNumber:      22,
				ConnectRetry:          7,
				AutoPosition:          true,
				SourceSsl:             sourceSSL,
			}
			node.ReplicaController = &replicaStatusController{status: status}
			iter, err := (&BaseBuilder{}).buildShowReplicaStatus(ctx, node, nil)
			require.NoError(t, err)
			defer iter.Close(ctx)

			actual, err := iter.Next(ctx)
			require.NoError(t, err)
			require.Len(t, actual, len(schema))
			require.Equal(t, "Ignored", actual[sslAllowedIndex])
			require.Equal(t, expected, actual)
			rows = append(rows, actual)
			_, err = iter.Next(ctx)
			require.ErrorIs(t, err, io.EOF)
		})
	}
	require.Equal(t, rows[0], rows[1])
}
