import { listSysConfigs, updateSysConfigs } from '@/services/sysconfig'

import { useConfigCenter } from './use-config-center'

export function useSysConfig() {
  return useConfigCenter(listSysConfigs, updateSysConfigs)
}
