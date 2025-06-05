from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.metrics_store_prometheus import MetricsStorePrometheus
    from ..models.metrics_store_victoria_metrics import MetricsStoreVictoriaMetrics


T = TypeVar("T", bound="MetricsStoreSpec")


@_attrs_define
class MetricsStoreSpec:
    """
    Attributes:
        prometheus (Union[Unset, MetricsStorePrometheus]):
        victoriametrics (Union[Unset, MetricsStoreVictoriaMetrics]):
    """

    prometheus: Union[Unset, "MetricsStorePrometheus"] = UNSET
    victoriametrics: Union[Unset, "MetricsStoreVictoriaMetrics"] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        prometheus: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.prometheus, Unset):
            prometheus = self.prometheus.to_dict()

        victoriametrics: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.victoriametrics, Unset):
            victoriametrics = self.victoriametrics.to_dict()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if prometheus is not UNSET:
            field_dict["prometheus"] = prometheus
        if victoriametrics is not UNSET:
            field_dict["victoriametrics"] = victoriametrics

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.metrics_store_prometheus import MetricsStorePrometheus
        from ..models.metrics_store_victoria_metrics import MetricsStoreVictoriaMetrics

        d = dict(src_dict)
        _prometheus = d.pop("prometheus", UNSET)
        prometheus: Union[Unset, MetricsStorePrometheus]
        if isinstance(_prometheus, Unset):
            prometheus = UNSET
        else:
            prometheus = MetricsStorePrometheus.from_dict(_prometheus)

        _victoriametrics = d.pop("victoriametrics", UNSET)
        victoriametrics: Union[Unset, MetricsStoreVictoriaMetrics]
        if isinstance(_victoriametrics, Unset):
            victoriametrics = UNSET
        else:
            victoriametrics = MetricsStoreVictoriaMetrics.from_dict(_victoriametrics)

        metrics_store_spec = cls(
            prometheus=prometheus,
            victoriametrics=victoriametrics,
        )

        metrics_store_spec.additional_properties = d
        return metrics_store_spec

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
