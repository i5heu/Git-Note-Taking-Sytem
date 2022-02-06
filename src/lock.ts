export interface LockStore {
    name: string,
    timeCreated: number
    timeOut: number //timeout is to control if a task should still run
}

// Handles the lock status of ongoing plugins, ongoing commits and ongoing pulls
export default class Lock {
    lockStore: LockStore[] = [];

    constructor() { }

    lock(name: string, timeOut: number) {
        //check if lock is already in store
        const locked = this.lockStore.find(l => l.name === name);
        if(locked) throw new Error("Lock already exists");

        this.lockStore.push({
            name,
            timeCreated: Date.now(),
            timeOut,
        });
    }

    unlock(name: string) {
        this.lockStore = this.lockStore.filter(l => l.name !== name);
    }

    listLocks() {
        return this.lockStore;
    }

    isLocked() {
        return this.lockStore.length > 0;
    }

    // will run a function when all locks are gone and timeOut is not expired
    async waitForFreeLockAndLock(name: LockStore["name"], timeOut: LockStore["timeOut"], callback: () => Promise<void|false>) {
        const timeOutUnix = Date.now() + timeOut * 1000;

        while (this.isLocked()) {
            await new Promise(resolve => setTimeout(resolve, 666));
        }

        if (Date.now() > timeOutUnix) return false;

        this.lock(name, timeOut);
        await callback();
        this.unlock(name);
    }
}