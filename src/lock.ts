import GitManager from "./gitManager";

export interface LockStore {
  name: string;
  timeCreated: number;
  timeOut: number; //timeout is to control if a task should still run
}

// Handles the lock status of ongoing plugins, ongoing commits and ongoing pulls
export default class Lock {
  lockStore: LockStore[] = [];
  queue: LockStore[] = [];
  emptyQueueFunction: GitManager["commitAndPush"];

  constructor() {}

  lock(name: string, timeOut: number) {
    //check if lock is already in store
    const locked = this.lockStore.find((l) => l.name === name);
    if (locked) throw new Error("Lock already exists :" + name);

    this.lockStore.push({
      name,
      timeCreated: Date.now(),
      timeOut,
    });
  }

  unlock(name: string, noUnlockCommit: boolean = false) {
    this.lockStore = this.lockStore.filter((l) => l.name !== name);
    this.queue = this.queue.filter((l) => l.name !== name);

    if (!noUnlockCommit) this.runEmptyQueueFunction();
  }

  listLocks() {
    return this.lockStore;
  }

  isLocked() {
    return this.lockStore.length > 0;
  }

  // will run a function when all locks are gone and timeOut is not expired
  async waitForFreeLockAndLock(
    name: LockStore["name"],
    timeOut: LockStore["timeOut"],
    callback: () => Promise<void | false>,
    noUnlockCommit: boolean = false
  ) {
    const timeOutUnix = Date.now() + timeOut * 1000;

    if (!noUnlockCommit)
      this.queue.push({
        name,
        timeCreated: Date.now(),
        timeOut,
      });

    while (this.isLocked()) {
      await new Promise((resolve) => setTimeout(resolve, 666));
    }

    if (Date.now() > timeOutUnix) {
      this.unlock(name, noUnlockCommit);
      return false;
    }

    this.lock(name, timeOut);
    await callback();
    this.unlock(name, noUnlockCommit);
  }

  queueLength() {
    return this.queue.length;
  }

  // run commit if last queue item is unlocked
  async runEmptyQueueFunction() {
    if (this.queue.length !== 0) return;
    this.emptyQueueFunction("Tyche: Empty Queue", true);
  }
}
