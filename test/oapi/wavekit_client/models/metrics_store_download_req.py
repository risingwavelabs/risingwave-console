import datetime
from collections.abc import Mapping
from typing import Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field
from dateutil.parser import isoparse

from ..types import UNSET, Unset

T = TypeVar("T", bound="MetricsStoreDownloadReq")


@_attrs_define
class MetricsStoreDownloadReq:
    """
    Attributes:
        step (Union[Unset, str]): Step of the metrics store, e.g. 1h, 1d, 1w, 1m, 1s
        start (Union[Unset, datetime.datetime]): Start time of the metrics store
        end (Union[Unset, datetime.datetime]): End time of the metrics store
        query_ratio (Union[Unset, float]): (0, 1], if OOM, reduce the memory usage in Prometheus instance by this ratio
            (default: 1)
        query (Union[Unset, str]): query to get the metrics, e.g. `{namespace="wavekit"}`
    """

    step: Union[Unset, str] = UNSET
    start: Union[Unset, datetime.datetime] = UNSET
    end: Union[Unset, datetime.datetime] = UNSET
    query_ratio: Union[Unset, float] = UNSET
    query: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        step = self.step

        start: Union[Unset, str] = UNSET
        if not isinstance(self.start, Unset):
            start = self.start.isoformat()

        end: Union[Unset, str] = UNSET
        if not isinstance(self.end, Unset):
            end = self.end.isoformat()

        query_ratio = self.query_ratio

        query = self.query

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if step is not UNSET:
            field_dict["step"] = step
        if start is not UNSET:
            field_dict["start"] = start
        if end is not UNSET:
            field_dict["end"] = end
        if query_ratio is not UNSET:
            field_dict["queryRatio"] = query_ratio
        if query is not UNSET:
            field_dict["query"] = query

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        step = d.pop("step", UNSET)

        _start = d.pop("start", UNSET)
        start: Union[Unset, datetime.datetime]
        if isinstance(_start, Unset):
            start = UNSET
        else:
            start = isoparse(_start)

        _end = d.pop("end", UNSET)
        end: Union[Unset, datetime.datetime]
        if isinstance(_end, Unset):
            end = UNSET
        else:
            end = isoparse(_end)

        query_ratio = d.pop("queryRatio", UNSET)

        query = d.pop("query", UNSET)

        metrics_store_download_req = cls(
            step=step,
            start=start,
            end=end,
            query_ratio=query_ratio,
            query=query,
        )

        metrics_store_download_req.additional_properties = d
        return metrics_store_download_req

    @property
    def additional_keys(self) -> list[str]:
        return list(self.additional_properties.keys())

    def __getitem__(self, key: str) -> Any:
        return self.additional_properties[key]

    def __setitem__(self, key: str, value: Any) -> None:
        self.additional_properties[key] = value

    def __delitem__(self, key: str) -> None:
        del self.additional_properties[key]

    def __contains__(self, key: str) -> bool:
        return key in self.additional_properties
