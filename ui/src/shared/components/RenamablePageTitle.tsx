// Libraries
import classnames from 'classnames'
import React, {
  ChangeEvent,
  FunctionComponent,
  MouseEvent,
  useEffect,
  useState,
  KeyboardEvent,
} from 'react'

// Components
import {Icon, IconFont, Input, InputRef, Page} from '@influxdata/clockface'
import {ClickOutside} from 'src/shared/components/ClickOutside'

interface Props {
  onRename: (name: string) => void
  onClickOutside?: (ev: MouseEvent<HTMLElement>) => void
  name: string
  placeholder: string
  maxLength: number
}

const RenamablePageTitle: FunctionComponent<Props> = ({
  onRename,
  onClickOutside,
  name,
  placeholder,
  maxLength,
}) => {
  const [isEditing, setEditing] = useState(false)
  const [workingName, setWorkingName] = useState(name)

  useEffect(() => {
    setWorkingName(name)

    return () => {
      setEditing(false)
    }
  }, [name])

  const handleStartEditing = (): void => {
    setEditing(true)
  }

  const handleStopEditing = (ev: MouseEvent<any>) => {
    onRename(workingName)

    if (onClickOutside) {
      onClickOutside(ev)
    }
  }

  const handleOnBlur = () => {
    setEditing(false)
  }

  const handleInputFocus = (ev: ChangeEvent<InputRef>) => {
    ev.currentTarget.select()
  }

  const handleInputChange = (ev: ChangeEvent<InputRef>) => {
    setWorkingName(ev.target.value)
  }

  const handleKeyDown = (ev: KeyboardEvent<InputRef>): void => {
    if (ev.key === 'Enter' || ev.key === 'Tab') {
      onRename(workingName)
      setEditing(false)
    }

    if (ev.key === 'Escape') {
      setEditing(false)
      setWorkingName(name)
    }
  }

  const renamablePageTitleClass = classnames('renamable-page-title', {
    untitled: name === placeholder || name === '',
  })

  if (isEditing) {
    return (
      <ClickOutside onClickOutside={handleStopEditing}>
        <div
          className={renamablePageTitleClass}
          data-testid="renamable-page-title"
        >
          <Input
            maxLength={maxLength}
            autoFocus={true}
            spellCheck={true}
            placeholder={placeholder}
            onBlur={handleOnBlur}
            // @ts-ignore
            onFocus={handleInputFocus}
            onChange={handleInputChange}
            onKeyDown={handleKeyDown}
            className={'renamable-page-title--input'}
            value={workingName}
            testID={'renamable-page-title--input'}
          />
        </div>
      </ClickOutside>
    )
  }

  return (
    <div className={renamablePageTitleClass} onClick={handleStartEditing}>
      <Page.Title title={workingName || placeholder} />
      <Icon glyph={IconFont.Pencil} className={'renamable-page-title--icon'} />
    </div>
  )
}

export default RenamablePageTitle
