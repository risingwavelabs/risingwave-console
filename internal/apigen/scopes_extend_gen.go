package apigen

import "github.com/gofiber/fiber/v2"

type Validator interface { 
    PreValidate(*fiber.Ctx) error
    
    PostValidate(*fiber.Ctx) error

    OwnDatabase(c *fiber.Ctx, UserID int32, DatabaseID int32) error
 
    GetOrgID(c *fiber.Ctx) int32
}


type XMiddleware struct {
	Handler ServerInterface
	Validator
}

func NewXMiddleware(handler ServerInterface, validator Validator) ServerInterface {
	return &XMiddleware{Handler: handler, Validator: validator}
}

// Refresh access token
// (POST /auth/refresh)
func (x *XMiddleware) RefreshToken(c *fiber.Ctx) error {
    
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.RefreshToken(c)
}
// Sign in user
// (POST /auth/sign-in)
func (x *XMiddleware) SignIn(c *fiber.Ctx) error {
    
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.SignIn(c)
}
// Sign out user
// (POST /auth/sign-out)
func (x *XMiddleware) SignOut(c *fiber.Ctx) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.SignOut(c)
}
// List all cluster versions
// (GET /cluster-versions)
func (x *XMiddleware) ListClusterVersions(c *fiber.Ctx) error {
    
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.ListClusterVersions(c)
}
// List all clusters
// (GET /clusters)
func (x *XMiddleware) ListClusters(c *fiber.Ctx) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.ListClusters(c)
}
// Import a cluster
// (POST /clusters/import)
func (x *XMiddleware) ImportCluster(c *fiber.Ctx) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.ImportCluster(c)
}
// Delete cluster
// (DELETE /clusters/{ID})
func (x *XMiddleware) DeleteCluster(c *fiber.Ctx, id int32, params DeleteClusterParams) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.DeleteCluster(c, id, params)
}
// Get cluster details
// (GET /clusters/{ID})
func (x *XMiddleware) GetCluster(c *fiber.Ctx, id int32) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.GetCluster(c, id)
}
// Update cluster
// (PUT /clusters/{ID})
func (x *XMiddleware) UpdateCluster(c *fiber.Ctx, id int32) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.UpdateCluster(c, id)
}
// Get snapshot configuration
// (GET /clusters/{ID}/auto-backup-config)
func (x *XMiddleware) GetClusterAutoBackupConfig(c *fiber.Ctx, id int32) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.GetClusterAutoBackupConfig(c, id)
}
// Update snapshot configuration
// (PUT /clusters/{ID}/auto-backup-config)
func (x *XMiddleware) UpdateClusterAutoBackupConfig(c *fiber.Ctx, id int32) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.UpdateClusterAutoBackupConfig(c, id)
}
// List diagnostic data
// (GET /clusters/{ID}/diagnostics)
func (x *XMiddleware) ListClusterDiagnostics(c *fiber.Ctx, id int32, params ListClusterDiagnosticsParams) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.ListClusterDiagnostics(c, id, params)
}
// Create diagnostic data
// (POST /clusters/{ID}/diagnostics)
func (x *XMiddleware) CreateClusterDiagnostic(c *fiber.Ctx, id int32) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.CreateClusterDiagnostic(c, id)
}
// Get diagnostic configuration
// (GET /clusters/{ID}/diagnostics/config)
func (x *XMiddleware) GetClusterAutoDiagnosticConfig(c *fiber.Ctx, id int32) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.GetClusterAutoDiagnosticConfig(c, id)
}
// Update diagnostic configuration
// (PUT /clusters/{ID}/diagnostics/config)
func (x *XMiddleware) UpdateClusterAutoDiagnosticConfig(c *fiber.Ctx, id int32) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.UpdateClusterAutoDiagnosticConfig(c, id)
}
// Get diagnostic data
// (GET /clusters/{ID}/diagnostics/{diagnosticId})
func (x *XMiddleware) GetClusterDiagnostic(c *fiber.Ctx, id int32, diagnosticId int32) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.GetClusterDiagnostic(c, id, diagnosticId)
}
// Run risectl command
// (POST /clusters/{ID}/risectl)
func (x *XMiddleware) RunRisectlCommand(c *fiber.Ctx, id int32) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.RunRisectlCommand(c, id)
}
// List cluster snapshots
// (GET /clusters/{ID}/snapshots)
func (x *XMiddleware) ListClusterSnapshots(c *fiber.Ctx, id int32) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.ListClusterSnapshots(c, id)
}
// Create a new snapshot
// (POST /clusters/{ID}/snapshots)
func (x *XMiddleware) CreateClusterSnapshot(c *fiber.Ctx, id int32) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.CreateClusterSnapshot(c, id)
}
// Delete snapshot
// (DELETE /clusters/{ID}/snapshots/{snapshotId})
func (x *XMiddleware) DeleteClusterSnapshot(c *fiber.Ctx, id int32, snapshotId int64) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.DeleteClusterSnapshot(c, id, snapshotId)
}
// Restore snapshot
// (POST /clusters/{ID}/snapshots/{snapshotId})
func (x *XMiddleware) RestoreClusterSnapshot(c *fiber.Ctx, id int32, snapshotId int64) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.RestoreClusterSnapshot(c, id, snapshotId)
}
// List all databases
// (GET /databases)
func (x *XMiddleware) ListDatabases(c *fiber.Ctx) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.ListDatabases(c)
}
// Import a database
// (POST /databases/import)
func (x *XMiddleware) ImportDatabase(c *fiber.Ctx) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.ImportDatabase(c)
}
// Test database connection
// (POST /databases/test-connection)
func (x *XMiddleware) TestDatabaseConnection(c *fiber.Ctx) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.TestDatabaseConnection(c)
}
// Delete database
// (DELETE /databases/{ID})
func (x *XMiddleware) DeleteDatabase(c *fiber.Ctx, id int32) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.DeleteDatabase(c, id)
}
// Get database details
// (GET /databases/{ID})
func (x *XMiddleware) GetDatabase(c *fiber.Ctx, id int32) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	 
	if err := x.OwnDatabase(c, x.GetOrgID(c), id); err != nil {
	    return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.GetDatabase(c, id)
}
// Update database
// (PUT /databases/{ID})
func (x *XMiddleware) UpdateDatabase(c *fiber.Ctx, id int32) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.UpdateDatabase(c, id)
}
// Get DDL progress
// (GET /databases/{ID}/ddl-progress)
func (x *XMiddleware) GetDDLProgress(c *fiber.Ctx, id int32) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.GetDDLProgress(c, id)
}
// Cancel DDL progress
// (POST /databases/{ID}/ddl-progress/{ddlID}/cancel)
func (x *XMiddleware) CancelDDLProgress(c *fiber.Ctx, id int32, ddlID int64) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.CancelDDLProgress(c, id, ddlID)
}
// Query database
// (POST /databases/{ID}/query)
func (x *XMiddleware) QueryDatabase(c *fiber.Ctx, id int32) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.QueryDatabase(c, id)
}
// Get all events
// (GET /events)
func (x *XMiddleware) ListEvents(c *fiber.Ctx) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.ListEvents(c)
}
// Get all metrics stores
// (GET /metrics-stores)
func (x *XMiddleware) ListMetricsStores(c *fiber.Ctx) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.ListMetricsStores(c)
}
// Import a metrics store
// (POST /metrics-stores/import)
func (x *XMiddleware) ImportMetricsStore(c *fiber.Ctx) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.ImportMetricsStore(c)
}
// Delete a metrics store
// (DELETE /metrics-stores/{ID})
func (x *XMiddleware) DeleteMetricsStore(c *fiber.Ctx, id int32, params DeleteMetricsStoreParams) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.DeleteMetricsStore(c, id, params)
}
// Get a metrics store
// (GET /metrics-stores/{ID})
func (x *XMiddleware) GetMetricsStore(c *fiber.Ctx, id int32) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.GetMetricsStore(c, id)
}
// Update a metrics store
// (PUT /metrics-stores/{ID})
func (x *XMiddleware) UpdateMetricsStore(c *fiber.Ctx, id int32) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.UpdateMetricsStore(c, id)
}
// Get materialized view throughput
// (GET /metrics/{clusterID}/materialized-view-throughput)
func (x *XMiddleware) GetMaterializedViewThroughput(c *fiber.Ctx, clusterID int32) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.GetMaterializedViewThroughput(c, clusterID)
}
// Get all tasks
// (GET /tasks)
func (x *XMiddleware) ListTasks(c *fiber.Ctx) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.ListTasks(c)
}
// Test cluster connection
// (POST /test-cluster-connection)
func (x *XMiddleware) TestClusterConnection(c *fiber.Ctx) error {
    if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	} 
	if err := x.PreValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	
	  
	if err := x.PostValidate(c); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
    return x.Handler.TestClusterConnection(c)
}

