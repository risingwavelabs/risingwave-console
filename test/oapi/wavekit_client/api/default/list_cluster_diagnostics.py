import datetime
from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ... import errors
from ...client import AuthenticatedClient, Client
from ...models.diagnostic_data import DiagnosticData
from ...types import UNSET, Response, Unset


def _get_kwargs(
    id: int,
    *,
    from_: Union[Unset, datetime.datetime] = UNSET,
    to: Union[Unset, datetime.datetime] = UNSET,
    page: Union[Unset, int] = 1,
    per_page: Union[Unset, int] = 20,
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    json_from_: Union[Unset, str] = UNSET
    if not isinstance(from_, Unset):
        json_from_ = from_.isoformat()
    params["from"] = json_from_

    json_to: Union[Unset, str] = UNSET
    if not isinstance(to, Unset):
        json_to = to.isoformat()
    params["to"] = json_to

    params["page"] = page

    params["perPage"] = per_page

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": f"/clusters/{id}/diagnostics",
        "params": params,
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Optional[list["DiagnosticData"]]:
    if response.status_code == 200:
        response_200 = []
        _response_200 = response.json()
        for response_200_item_data in _response_200:
            response_200_item = DiagnosticData.from_dict(response_200_item_data)

            response_200.append(response_200_item)

        return response_200
    if client.raise_on_unexpected_status:
        raise errors.UnexpectedStatus(response.status_code, response.content)
    else:
        return None


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[list["DiagnosticData"]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    id: int,
    *,
    client: AuthenticatedClient,
    from_: Union[Unset, datetime.datetime] = UNSET,
    to: Union[Unset, datetime.datetime] = UNSET,
    page: Union[Unset, int] = 1,
    per_page: Union[Unset, int] = 20,
) -> Response[list["DiagnosticData"]]:
    """List diagnostic data

     Retrieve diagnostic data for a specific cluster with optional date range filtering

    Args:
        id (int):
        from_ (Union[Unset, datetime.datetime]):
        to (Union[Unset, datetime.datetime]):
        page (Union[Unset, int]):  Default: 1.
        per_page (Union[Unset, int]):  Default: 20.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[list['DiagnosticData']]
    """

    kwargs = _get_kwargs(
        id=id,
        from_=from_,
        to=to,
        page=page,
        per_page=per_page,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    id: int,
    *,
    client: AuthenticatedClient,
    from_: Union[Unset, datetime.datetime] = UNSET,
    to: Union[Unset, datetime.datetime] = UNSET,
    page: Union[Unset, int] = 1,
    per_page: Union[Unset, int] = 20,
) -> Optional[list["DiagnosticData"]]:
    """List diagnostic data

     Retrieve diagnostic data for a specific cluster with optional date range filtering

    Args:
        id (int):
        from_ (Union[Unset, datetime.datetime]):
        to (Union[Unset, datetime.datetime]):
        page (Union[Unset, int]):  Default: 1.
        per_page (Union[Unset, int]):  Default: 20.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        list['DiagnosticData']
    """

    return sync_detailed(
        id=id,
        client=client,
        from_=from_,
        to=to,
        page=page,
        per_page=per_page,
    ).parsed


async def asyncio_detailed(
    id: int,
    *,
    client: AuthenticatedClient,
    from_: Union[Unset, datetime.datetime] = UNSET,
    to: Union[Unset, datetime.datetime] = UNSET,
    page: Union[Unset, int] = 1,
    per_page: Union[Unset, int] = 20,
) -> Response[list["DiagnosticData"]]:
    """List diagnostic data

     Retrieve diagnostic data for a specific cluster with optional date range filtering

    Args:
        id (int):
        from_ (Union[Unset, datetime.datetime]):
        to (Union[Unset, datetime.datetime]):
        page (Union[Unset, int]):  Default: 1.
        per_page (Union[Unset, int]):  Default: 20.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[list['DiagnosticData']]
    """

    kwargs = _get_kwargs(
        id=id,
        from_=from_,
        to=to,
        page=page,
        per_page=per_page,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    id: int,
    *,
    client: AuthenticatedClient,
    from_: Union[Unset, datetime.datetime] = UNSET,
    to: Union[Unset, datetime.datetime] = UNSET,
    page: Union[Unset, int] = 1,
    per_page: Union[Unset, int] = 20,
) -> Optional[list["DiagnosticData"]]:
    """List diagnostic data

     Retrieve diagnostic data for a specific cluster with optional date range filtering

    Args:
        id (int):
        from_ (Union[Unset, datetime.datetime]):
        to (Union[Unset, datetime.datetime]):
        page (Union[Unset, int]):  Default: 1.
        per_page (Union[Unset, int]):  Default: 20.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        list['DiagnosticData']
    """

    return (
        await asyncio_detailed(
            id=id,
            client=client,
            from_=from_,
            to=to,
            page=page,
            per_page=per_page,
        )
    ).parsed
