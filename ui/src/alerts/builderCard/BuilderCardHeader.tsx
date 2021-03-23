// Libraries
import React from 'react'
import classnames from 'classnames'

interface Props {
  title: string
  testID?: string
  onDelete?: () => void
  onDragStart?: () => void
  className?: string
}

const BuilderCardHeader: React.FC<Props> = props => {
  const {testID = 'builder-card--header', className, children, title} = props

  const classname = classnames('builder-card--header', {
    [`${className}`]: className,
  })

  return (
    <div className={classname} data-testid={testID}>
      <h2 className={'builder-card--title'}>{title}</h2>
      {children}
    </div>
  )
}

export default BuilderCardHeader
