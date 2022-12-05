// Libraries
import React, {FunctionComponent, RefObject, useRef, useState} from 'react'
import classnames from 'classnames'

// Components
import ContextItem from 'src/cells/components/ContextItem'
import {
  Appearance,
  Icon,
  IconFont,
  Popover,
  PopoverInteraction,
} from '@influxdata/clockface'
import ContextDangerItem from 'src/cells/components/ContextDangerItem'

// Types
import {Cell} from 'src/types/cells'

// Hooks
import {useNavigate} from 'react-router-dom'
import {useDispatch, useSelector} from 'react-redux'

// Actions
import {removeCell} from 'src/cells/actions/thunk'
import {getDashboardID} from 'src/dashboards/selectors'

interface Props {
  cell: Cell
}

const Context: FunctionComponent<Props> = ({cell}) => {
  const dispatch = useDispatch()
  const dashboardID = useSelector(getDashboardID)
  const navigate = useNavigate()

  const handleEditCell = (): void => {
    navigate(`${window.location.pathname}/cells/${cell.id}/edit`)
  }

  const handleEditNote = () => {
    console.log('edit note')
  }

  const handleDeleteCell = () => {
    dispatch(removeCell(dashboardID, cell.id))
  }

  const popoverContents = (onHide?: () => void): JSX.Element => {
    return (
      <div className={'cell--context-menu'}>
        <ContextItem
          label={'Configure'}
          icon={IconFont.Pencil}
          onClick={handleEditCell}
          onHide={onHide}
          testID={'cell-context--configure'}
        />

        <ContextItem
          label={'Add Note'}
          icon={IconFont.Text_New}
          onClick={handleEditNote}
          onHide={onHide}
          testID={'cell-context--note'}
        />

        <ContextDangerItem
          label={'Delete'}
          onClick={handleDeleteCell}
          icon={IconFont.Trash_New}
          onHide={onHide}
          testID={'cell-context--delete'}
        />
      </div>
    )
  }

  const [popoverVisible, setPopoverVisibility] = useState(false)
  const buttonClass = classnames('cell--context', {
    'cell--context__active': popoverVisible,
  })
  const triggerRef: RefObject<HTMLButtonElement> =
    useRef<HTMLButtonElement>(null)

  return (
    <>
      <button className={buttonClass} ref={triggerRef}>
        <Icon glyph={IconFont.CogOutline_New} />
      </button>

      <Popover
        appearance={Appearance.Outline}
        enableDefaultStyles={false}
        showEvent={PopoverInteraction.Click}
        hideEvent={PopoverInteraction.Click}
        triggerRef={triggerRef}
        contents={popoverContents}
        onShow={() => setPopoverVisibility(true)}
        onHide={() => {
          setPopoverVisibility(false)
        }}
      />
    </>
  )
}

export default Context
