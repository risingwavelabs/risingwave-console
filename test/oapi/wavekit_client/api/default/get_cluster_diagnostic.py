from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ... import errors
from ...client import AuthenticatedClient, Client
from ...models.diagnostic_data import DiagnosticData
from ...types import Response


def _get_kwargs(
    id: int,
    diagnostic_id: int,
) -> dict[str, Any]:
    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": f"/clusters/{id}/diagnostics/{diagnostic_id}",
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Optional[DiagnosticData]:
    if response.status_code == 200:
        response_200 = DiagnosticData.from_dict(response.json())

        return response_200
    if client.raise_on_unexpected_status:
        raise errors.UnexpectedStatus(response.status_code, response.content)
    else:
        return None


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[DiagnosticData]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    id: int,
    diagnostic_id: int,
    *,
    client: AuthenticatedClient,
) -> Response[DiagnosticData]:
    """Get diagnostic data

     Get diagnostic data for a specific cluster

    Args:
        id (int):
        diagnostic_id (int):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[DiagnosticData]
    """

    kwargs = _get_kwargs(
        id=id,
        diagnostic_id=diagnostic_id,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    id: int,
    diagnostic_id: int,
    *,
    client: AuthenticatedClient,
) -> Optional[DiagnosticData]:
    """Get diagnostic data

     Get diagnostic data for a specific cluster

    Args:
        id (int):
        diagnostic_id (int):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        DiagnosticData
    """

    return sync_detailed(
        id=id,
        diagnostic_id=diagnostic_id,
        client=client,
    ).parsed


async def asyncio_detailed(
    id: int,
    diagnostic_id: int,
    *,
    client: AuthenticatedClient,
) -> Response[DiagnosticData]:
    """Get diagnostic data

     Get diagnostic data for a specific cluster

    Args:
        id (int):
        diagnostic_id (int):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[DiagnosticData]
    """

    kwargs = _get_kwargs(
        id=id,
        diagnostic_id=diagnostic_id,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    id: int,
    diagnostic_id: int,
    *,
    client: AuthenticatedClient,
) -> Optional[DiagnosticData]:
    """Get diagnostic data

     Get diagnostic data for a specific cluster

    Args:
        id (int):
        diagnostic_id (int):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        DiagnosticData
    """

    return (
        await asyncio_detailed(
            id=id,
            diagnostic_id=diagnostic_id,
            client=client,
        )
    ).parsed
