import React from "react";
import { Icon, IconFont } from "@influxdata/clockface";

interface Props {
  testID?: string
  label: string
  icon: IconFont
  onClick: () => void
  onHide?: () => void
}

const CellContextItem: React.FC<Props> = ({
  icon,
  label,
  testID,
  onHide,
  onClick
}) => {
  const handleClick = (): void => {
    onHide && onHide();
    onClick();
  };

  return (
    <div
      className="cell--context-item"
      onClick={handleClick}
      data-testid={testID}
    >
      <Icon glyph={icon} />
      {label}
    </div>
  );
};

export default CellContextItem;
