import React from 'react'

const VERSION = process.env.VERSION
const GIT_SHA = process.env.GIT_SHA

const VersionInfo: React.FC = props => {
  return (
    <div className={'version-info'}>
      <p>
        Version {VERSION} {GIT_SHA && <code>({GIT_SHA.slice(0, 7)}</code>}
      </p>
    </div>
  )
}

export default VersionInfo
