import React, {FC} from 'react'

interface Props {
  children: any
}

export const WizardProgressHeader: FC<Props> = props => {
  const {children} = props

  return <div className={'wizard--progress-header'}>{children}</div>
}
