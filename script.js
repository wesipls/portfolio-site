function prepareAndPrint() {
    const originalTitle = document.title;
    // This sets the filename for the "Save as PDF" dialog
    document.title = "Kim_Valkama_Infrastructure_and_Platform_Engineer_CV"; 
    
    window.print();
    
    // Restore the original browser tab title
    document.title = originalTitle;
}
console.log(
    `%cHope you like Vim; if we work together, that's going to be your new default server editor.`, 
    "color: #4af626; font-family: monospace; font-size: 13px; line-height: 1.5;"
);
