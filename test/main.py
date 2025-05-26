import pytest
import httpx

from oapi.wavekit_client.client import AuthenticatedClient

base_url = "http://localhost:8020/api/v1"


@pytest.fixture
def auth_client():
    response = httpx.post(
        f"{base_url}/auth/sign-in", data={"name": "root", "password": "123456"}
    )
    assert response.status_code == 200
    access_token = response.json()["accessToken"]
    return AuthenticatedClient(base_url=base_url, token=access_token)


def test_cluster_db_management(auth_client: AuthenticatedClient):
    from oapi.wavekit_client.api.default import (
        list_clusters,
        get_cluster,
        import_cluster,
        delete_cluster,
    )
    from oapi.wavekit_client.models import ClusterImport

    # list clusters
    response = list_clusters.sync_detailed(client=auth_client)
    assert response.status_code == 200
    assert len(response.parsed) > 0

    # get cluster
    response = get_cluster.sync_detailed(client=auth_client, id=response.parsed[0].id)
    assert response.status_code == 200
    assert response.parsed is not None

    # import cluster
    response = import_cluster.sync_detailed(
        client=auth_client,
        body=ClusterImport(
            name="test",
            host="localhost",
            sql_port=5432,
            meta_port=9191,
            http_port=8080,
            version="v1.0.0",
            metrics_store_id=1,
        ),
    )

    imported_cluster = response.parsed
    assert response.status_code == 201
    assert imported_cluster is not None

    # get the imported cluster
    response = get_cluster.sync_detailed(client=auth_client, id=imported_cluster.id)
    assert response.status_code == 200
    assert response.parsed.id == imported_cluster.id
    assert response.parsed.name == imported_cluster.name
    assert response.parsed.host == imported_cluster.host
    assert response.parsed.sql_port == imported_cluster.sql_port
    assert response.parsed.meta_port == imported_cluster.meta_port
    assert response.parsed.http_port == imported_cluster.http_port
    assert response.parsed.version == imported_cluster.version
    assert response.parsed.metrics_store_id == imported_cluster.metrics_store_id

    # delete cluster
    response = delete_cluster.sync_detailed(client=auth_client, id=imported_cluster.id)
    assert response.status_code == 204



if __name__ == "__main__":
    pytest.main([__file__])
