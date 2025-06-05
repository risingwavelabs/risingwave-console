from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..models.credentials_token_type import CredentialsTokenType

T = TypeVar("T", bound="Credentials")


@_attrs_define
class Credentials:
    """
    Attributes:
        access_token (str): JWT access token
        refresh_token (str): JWT refresh token for obtaining new access tokens
        token_type (CredentialsTokenType): Token type
    """

    access_token: str
    refresh_token: str
    token_type: CredentialsTokenType
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        access_token = self.access_token

        refresh_token = self.refresh_token

        token_type = self.token_type.value

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "accessToken": access_token,
                "refreshToken": refresh_token,
                "tokenType": token_type,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        access_token = d.pop("accessToken")

        refresh_token = d.pop("refreshToken")

        token_type = CredentialsTokenType(d.pop("tokenType"))

        credentials = cls(
            access_token=access_token,
            refresh_token=refresh_token,
            token_type=token_type,
        )

        credentials.additional_properties = d
        return credentials

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
