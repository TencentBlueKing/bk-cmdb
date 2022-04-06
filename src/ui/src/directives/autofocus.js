/**
 * @directive 自动聚焦指令，目前只支持 input
 */
export const autofocus = {
  update: (el) => {
    const input = el.querySelector('input')

    if (input) {
      // 尽量靠后执行，避免其他队列任务影响到聚焦
      setTimeout(() => {
        input.focus()
      }, 0)
    }
  }
}
