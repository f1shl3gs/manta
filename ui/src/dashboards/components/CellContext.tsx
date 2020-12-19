import React, { RefObject, useRef, useState } from "react";
import { Cell, ViewProperties } from "../../types";
import { Appearance, Icon, IconFont, Popover, PopoverInteraction } from "@influxdata/clockface";
import CellContextItem from "./CellContextItem";
import CellContextDangerItem from "./CellContextDangerItem";
import classnames from "classnames";

interface Props {
  cell: Cell
  view: ViewProperties
}

const CellContext: React.FC<Props> = props => {
  const {
    cell
  } = props;

  const handleEditCell = (): void => {
    console.log("edit cell");
  };

  const handleEditNote = () => {
    console.log("edit cell note");
  };

  const handleDeleteCell = () => {
    console.log("delete cell");
  };

  const editNoteText = cell.desc;

  const popoverContents = (onHide?: () => void): JSX.Element => {
    return (
      <div className="cell--context-menu">
        <CellContextItem
          label="Configure"
          onClick={handleEditCell}
          icon={IconFont.Pencil}
          onHide={onHide}
          testID="cell-context--configure"
        />
        <CellContextItem
          label={editNoteText}
          onClick={handleEditNote}
          icon={IconFont.TextBlock}
          onHide={onHide}
          testID="cell-context--note"
        />

        {/*
        <CellContextItem
          label="Clone"
          onClick={handleCloneCell}
          icon={IconFont.Duplicate}
          onHide={onHide}
          testID="cell-context--clone"
        />
*/}

        <CellContextDangerItem
          label="Delete"
          onClick={handleDeleteCell}
          icon={IconFont.Trash}
          onHide={onHide}
          testID="cell-context--delete"
        />
      </div>
    );
  };

  const [popoverVisible, setPopoverVisibility] = useState<boolean>(false);
  const buttonClass = classnames("cell--context", {
    "cell--context__active": popoverVisible
  });

  const triggerRef: RefObject<HTMLButtonElement> = useRef<HTMLButtonElement>(
    null
  );

  return (
    <>
      <button
        className={buttonClass}
        ref={triggerRef}
      >
        <Icon glyph={IconFont.CogThick} />
      </button>

      <Popover
        appearance={Appearance.Outline}
        enableDefaultStyles={false}
        showEvent={PopoverInteraction.Click}
        hideEvent={PopoverInteraction.Click}
        triggerRef={triggerRef}
        contents={popoverContents}
        onShow={() => {
          setPopoverVisibility(true);
        }}
        onHide={() => {
          setPopoverVisibility(false);
        }}
      />
    </>
  );
};

export default CellContext;
