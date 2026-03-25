export function prepareAndPrint() {
    const originalTitle = document.title;

    document.title = "Kim_Valkama_Infrastructure_and_Platform_Engineer_CV";
    window.print();

    document.title = originalTitle;
}
window.prepareAndPrint = prepareAndPrint;
