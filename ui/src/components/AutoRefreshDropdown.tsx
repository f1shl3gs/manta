import React, { useState } from "react";
import { ComponentStatus, Dropdown, IconFont } from "@influxdata/clockface";
import { AutoRefresh, AutoRefreshOption, AutoRefreshStatus } from "../types/AutoRefresh";

const dropdownIcon = (autoRefresh: AutoRefresh): IconFont => {
  if (autoRefresh.status) {
    return IconFont.Pause;
  }

  return IconFont.Refresh;
};

const dropdownStatus = (autoRefresh: AutoRefresh): ComponentStatus => {
  if (autoRefresh.status === AutoRefreshStatus.Disabled) {
    return ComponentStatus.Disabled;
  }

  return ComponentStatus.Default;
};

interface Props {
  options: AutoRefreshOption[]
}

const AutoRefreshDropdown: React.FC<Props> = props => {
  const {
    options
  } = props;

  const [autoRefresh, setAutoRefresh] = useState<AutoRefresh>({
    status: AutoRefreshStatus.Active,
    interval: 60
  });
  const [selected, setSelected] = useState(options[0]);

  return (
    <>
      <Dropdown
        button={(active, onClick) => (
          <Dropdown.Button
            active={active}
            onClick={onClick}
            status={dropdownStatus(autoRefresh)}
            icon={dropdownIcon(autoRefresh)}
          >
            {selected.label}
          </Dropdown.Button>
        )}
        menu={onCollapse => (
          <Dropdown.Menu
            onCollapse={onCollapse}
          >
            {
              options.map(option => (
                <Dropdown.Item
                  key={option.label}
                  value={option}
                  onClick={(v) => setSelected(v)}
                >
                  {option.label}
                </Dropdown.Item>
              ))
            }
          </Dropdown.Menu>
        )}
      />
    </>
  );
};

export default AutoRefreshDropdown;
