from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, cast

from attrs import define as _attrs_define
from attrs import field as _attrs_field

if TYPE_CHECKING:
    from ..models.metric_series_metric import MetricSeriesMetric


T = TypeVar("T", bound="MetricSeries")


@_attrs_define
class MetricSeries:
    """
    Attributes:
        metric (MetricSeriesMetric):
        values (list[list[Any]]):
    """

    metric: "MetricSeriesMetric"
    values: list[list[Any]]
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        metric = self.metric.to_dict()

        values = []
        for values_item_data in self.values:
            values_item = values_item_data

            values.append(values_item)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "metric": metric,
                "values": values,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.metric_series_metric import MetricSeriesMetric

        d = dict(src_dict)
        metric = MetricSeriesMetric.from_dict(d.pop("metric"))

        values = []
        _values = d.pop("values")
        for values_item_data in _values:
            values_item = cast(list[Any], values_item_data)

            values.append(values_item)

        metric_series = cls(
            metric=metric,
            values=values,
        )

        metric_series.additional_properties = d
        return metric_series

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
