import React, {FunctionComponent, RefObject, useRef, useState} from 'react'
import {Cell} from 'src/types/Dashboard'
import {useNavigate, useParams} from 'react-router-dom'
import ContextItem from './ContextItem'
import {
  Appearance,
  Icon,
  IconFont,
  Popover,
  PopoverInteraction,
} from '@influxdata/clockface'
import ContextDangerItem from './ContextDangerItem'
import classnames from 'classnames'
import useFetch from 'src/shared/useFetch'
import {
  defaultErrorNotification,
  useNotify,
} from 'src/shared/components/notifications/useNotification'
import {useDashboard} from '../../useDashboard'

interface Props {
  cell: Cell
}

const Context: FunctionComponent<Props> = ({cell}) => {
  const navigate = useNavigate()
  const notify = useNotify()
  const {dashboardId} = useParams()
  const {reload} = useDashboard()
  const {run: deleteCell} = useFetch(
    `/api/v1/dashboards/${dashboardId}/cells/${cell.id}`,
    {
      method: 'DELETE',
      onError: err => {
        notify({
          ...defaultErrorNotification,
          message: `Delete Cell ${cell.name} failed, ${err}`,
        })
      },
      onSuccess: _ => {
        reload()
      },
    }
  )

  const handleEditCell = (): void => {
    navigate(`${window.location.pathname}/cells/${cell.id}/edit`)
  }

  const handleEditNote = () => {
    console.log('edit note')
  }

  const handleDeleteCell = () => {
    deleteCell()
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
        onShow={() => {
          setPopoverVisibility(true)
        }}
        onHide={() => {
          setPopoverVisibility(false)
        }}
      />
    </>
  )
}

export default Context
