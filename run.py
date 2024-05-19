import os
import time
import signal
import subprocess


def execute_shell_prompt(prompt):
    try:
        print("Executing: {}".format(prompt))
        output = subprocess.check_output(prompt, shell=True, stderr=subprocess.STDOUT)
        return output.decode('utf-8')
    except subprocess.CalledProcessError as e:
        print(e.output.decode('utf-8'))
        raise e


def tail_logs(tail_command, log_file):
    process = None

    def quit_tail():
        print("Quitting tail")
        nonlocal process
        if process and process.poll() is None:
            os.killpg(os.getpgid(process.pid), signal.SIGTERM)
            # process.terminate()
            process = None

    try:
        cmd = tail_command + " > " + log_file
        process = subprocess.Popen(cmd, stdout=subprocess.PIPE, shell=True, preexec_fn=os.setsid)
        return quit_tail

    except subprocess.CalledProcessError as e:
        return e.output.decode('utf-8')

def workload():
    execute_shell_prompt('kubectl apply -f zookeeper-312/zkc-1.yaml')
    time.sleep(5)
    execute_shell_prompt('kubectl wait --for=condition=Ready pod/zookeeper-cluster-0 --timeout=60s -n default')
    execute_shell_prompt('kubectl delete -f zookeeper-312/zkc-1.yaml')
    execute_shell_prompt('kubectl wait --for=delete pod/zookeeper-cluster-0 --timeout=60s -n default')
    execute_shell_prompt('kubectl apply -f zookeeper-312/zkc-1.yaml')
    time.sleep(5)
    execute_shell_prompt('kubectl wait --for=condition=Ready pod/zookeeper-cluster-0 --timeout=60s -n default')


def main():
    print(execute_shell_prompt('ls -l'))

    for i in range(3):
        print("running iteration: {}".format(i))
        # tail and follow any future logs
        quit_tail = tail_logs('kubectl logs deployment/zookeeper-operator -n default --tail=0 -f', 'out-{}.json'.format(i))
        workload()
        print("done")
        quit_tail()


if __name__ == '__main__':
    main()
