from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field

T = TypeVar("T", bound="RisectlCommandResult")


@_attrs_define
class RisectlCommandResult:
    """
    Attributes:
        stdout (str): Standard output of the risectl command
        stderr (str): Standard error of the risectl command
        exit_code (int): Exit code of the risectl command
        err (str): Error message when try to run the risectl command
    """

    stdout: str
    stderr: str
    exit_code: int
    err: str
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        stdout = self.stdout

        stderr = self.stderr

        exit_code = self.exit_code

        err = self.err

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "stdout": stdout,
                "stderr": stderr,
                "exitCode": exit_code,
                "err": err,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        stdout = d.pop("stdout")

        stderr = d.pop("stderr")

        exit_code = d.pop("exitCode")

        err = d.pop("err")

        risectl_command_result = cls(
            stdout=stdout,
            stderr=stderr,
            exit_code=exit_code,
            err=err,
        )

        risectl_command_result.additional_properties = d
        return risectl_command_result

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
