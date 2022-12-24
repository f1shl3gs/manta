// Libraries
import classnames from 'classnames'
import React, {FunctionComponent} from 'react'

interface Props {
  title: string
  testID?: string
  className?: string
  children?: JSX.Element

  onDragStart?: () => void
}

const BuilderCardHeader: FunctionComponent<Props> = ({
  title,
  testID = 'builder-card-header',
  className,
  onDragStart,
  children,
}) => {
  const classname = classnames('builder-card--header', {
    [`${className}`]: className,
  })

  const titleContent = () => {
    if (onDragStart) {
      return (
        <div className={'builder-card--draggable'} onDragStart={onDragStart}>
          <div className={'builder-card--hamburger'} />
          <h2>{title}</h2>
        </div>
      )
    }

    return <h2 className={'builder-card--title'}>{title}</h2>
  }

  return (
    <div className={classname} data-testid={testID}>
      {titleContent()}
      {children}
    </div>
  )
}

export default BuilderCardHeader
