// Libraries
import React, { useCallback, useState } from 'react';
import classnames from 'classnames';

// Components
import {
  ComponentStatus,
  Dropdown,
  IconFont,
  SquareButton
} from '@influxdata/clockface';

// Hooks
import { useAutoRefresh } from 'shared/useAutoRefresh';

// Types
import {
  AutoRefresh,
  AutoRefreshOption,
  AutoRefreshOptionType,
  AutoRefreshStatus
} from 'types/AutoRefresh';

// Constants
import {
  AutoRefreshDropdownOptions,
  DROPDOWN_WIDTH_FULL,
  DROPDOWN_WIDTH_COLLAPSED
} from 'constants/autoRefresh';

const dropdownIcon = (autoRefresh: AutoRefresh): IconFont => {
  if (autoRefresh.status === AutoRefreshStatus.Paused) {
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

const AutoRefreshDropdown: React.FC<Props> = (props) => {
  const { options } = props;
  const { autoRefresh, setAutoRefresh } = useAutoRefresh();
  const [selected, setSelected] = useState(AutoRefreshDropdownOptions[3]);
  const paused = selected.seconds === 0;
  const dropdownWidthPixels = paused ? DROPDOWN_WIDTH_COLLAPSED : DROPDOWN_WIDTH_FULL;
  const dropdownClassname = classnames('autorefresh-dropdown', {
    paused: paused
  });
  const onSelectAutoRefreshOption = useCallback((opt: AutoRefreshOption) => {
    setSelected(opt);
    setAutoRefresh({
      status: opt.seconds !== 0 ? AutoRefreshStatus.Active : AutoRefreshStatus.Paused,
      interval: opt.seconds
    });
  }, []);

  return (
    <div className={dropdownClassname}>
      <Dropdown
        style={{ width: `${dropdownWidthPixels}px` }}
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
        menu={(onCollapse) => (
          <Dropdown.Menu
            onCollapse={onCollapse}
            style={{ width: `${DROPDOWN_WIDTH_FULL}px` }}
          >
            {options.map((option) => {
              if (option.type === AutoRefreshOptionType.Header) {
                return (
                  <Dropdown.Divider
                    key={option.id}
                    id={option.id}
                    text={option.label}
                  />
                );
              }

              return (
                <Dropdown.Item
                  key={option.id}
                  id={option.id}
                  value={option}
                  onClick={onSelectAutoRefreshOption}
                >
                  {option.label}
                </Dropdown.Item>
              );
            })}
          </Dropdown.Menu>
        )}
      />

      {!paused ? null : (
        <SquareButton
          icon={IconFont.Refresh}
          className={'autorefresh-dropdown--pause'}
        />
      )}
    </div>
  );
};

export default AutoRefreshDropdown;
