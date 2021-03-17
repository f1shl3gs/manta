// Libraries
import React from 'react'
import classnames from 'classnames'

interface Props {
  testID: string
  className?: string
  widthPixels?: number
}

const BuilderCard: React.FC<Props> = props => {
  const {testID, className, children, widthPixels = 228} = props

  const style = {flex: `0 0 ${widthPixels}px`}
  const classname = classnames('builder-card', {[`${className}`]: className})

  return (
    <div className={classname} data-testid={testID} style={style}>
      {children}
    </div>
  )
}

export default BuilderCard
