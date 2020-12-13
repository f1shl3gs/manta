import React, { ChangeEvent, useCallback } from "react";
import { IconFont, Input } from "@influxdata/clockface";

interface Props {
  onSearch: (searchTerm: string) => void
  placeholder: string
}

const SearchWidget: React.FC<Props> = props => {
  const {
    onSearch,
    placeholder,
  } = props;

  const onChange = useCallback((ev: ChangeEvent<HTMLInputElement>) => {
    onSearch(ev.target.value);
  }, [onSearch]);

  return (
    <Input
      icon={IconFont.Search}
      placeholder={placeholder}
      onChange={onChange}
      className="search-widget-input"
    />
  );
};

export default SearchWidget;
